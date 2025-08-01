package storage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olenka-91/DocsServer/internal/entity"
	"github.com/sirupsen/logrus"
)

const (
	cacheTTL             = 5 * time.Minute
	maxMemoryCacheSize   = 100 * 1024 * 1024 // 100MB
	maxMemoryCachedFiles = 100
)

type Cache struct {
	memoryCache *MemoryCache
}

type MemoryCache struct {
	sync.RWMutex
	files      map[uuid.UUID]*CachedFile
	totalSize  int64
	maxSize    int64
	maxEntries int
}

type CachedFile struct {
	data    []byte
	size    int64
	mime    string
	etag    string
	created time.Time
}

type FileStorage struct {
	basePath  string
	cache     *Cache
	fileLocks *sync.Map
}

func NewFileStorage(basePath string) *FileStorage {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		logrus.Fatalf("Failed to create storage directory: %v", err)
	}

	return &FileStorage{
		basePath: basePath,
		cache: &Cache{
			memoryCache: &MemoryCache{
				files:      make(map[uuid.UUID]*CachedFile),
				maxSize:    maxMemoryCacheSize,
				maxEntries: maxMemoryCachedFiles,
			},
		},
		fileLocks: &sync.Map{},
	}
}

func (fs *FileStorage) getFileLock(id uuid.UUID) *sync.Mutex {
	lock, _ := fs.fileLocks.LoadOrStore(id, &sync.Mutex{})
	return lock.(*sync.Mutex)
}

func (fs *FileStorage) getFilePath(id uuid.UUID, filename string) string {
	// Создаем поддиректории для распределения файлов
	subDir := id.String()[0:2]
	dirPath := filepath.Join(fs.basePath, subDir)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		logrus.Errorf("Failed to create subdirectory: %v", err)
		return filepath.Join(fs.basePath, id.String()+"_"+filename)
	}

	return filepath.Join(dirPath, id.String()+"_"+filename)
}

func (fs *FileStorage) SaveFile(id uuid.UUID, r io.Reader, filename string) (int64, string, error) {
	lock := fs.getFileLock(id)
	lock.Lock()
	defer lock.Unlock()

	filePath := fs.getFilePath(id, filename)
	tmpPath := filePath + ".tmp"

	// Создаем временный файл
	file, err := os.Create(tmpPath)
	if err != nil {
		return 0, "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Копируем данные и вычисляем хеш
	hasher := sha256.New()
	tee := io.TeeReader(r, hasher)

	size, err := io.Copy(file, tee)
	if err != nil {
		os.Remove(tmpPath)
		return 0, "", fmt.Errorf("failed to write file: %w", err)
	}

	// Финализируем запись
	if err := file.Sync(); err != nil {
		os.Remove(tmpPath)
		return 0, "", fmt.Errorf("failed to sync file: %w", err)
	}
	file.Close()

	// Переименовываем временный файл в постоянный
	if err := os.Rename(tmpPath, filePath); err != nil {
		os.Remove(tmpPath)
		return 0, "", fmt.Errorf("failed to rename file: %w", err)
	}

	// Определяем MIME-тип
	mimeType := mime.TypeByExtension(filepath.Ext(filename))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	//Обновляем кэш
	etag := hex.EncodeToString(hasher.Sum(nil))
	fs.cache.memoryCache.Store(id, &CachedFile{
		data:    nil, // Не кэшируем данные при сохранении
		size:    size,
		mime:    mimeType,
		etag:    etag,
		created: time.Now(),
	})

	return size, mimeType, nil
}

func (fs *FileStorage) DeleteFile(id uuid.UUID, filename string) error {
	lock := fs.getFileLock(id)
	lock.Lock()
	defer lock.Unlock()

	filePath := fs.getFilePath(id, filename)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Инвалидируем кэш
	fs.cache.memoryCache.Delete(id)
	return nil
}

func (fs *FileStorage) ServeFile(ctx *gin.Context, doc *entity.Document) error {
	// Для HEAD запросов только метаданные
	if ctx.Request.Method == http.MethodHead {
		return fs.serveFileHead(ctx.Writer, doc)
	}

	// Проверяем кэш
	if cached, ok := fs.cache.Get(doc.ID); ok {
		logrus.Debugf("Serving file %s from cache", doc.ID)
		ctx.Writer.Header().Set("Content-Type", cached.mime)
		ctx.Writer.Header().Set("ETag", cached.etag)
		ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", cached.size))
		http.ServeContent(ctx.Writer, ctx.Request, doc.Name, time.Now(), bytes.NewReader(cached.data))
		return nil
	}

	// Если не в кэше, читаем с диска
	return fs.serveFileFromDisk(ctx.Writer, ctx.Request, doc)
}

func (fs *FileStorage) serveFileHead(w http.ResponseWriter, doc *entity.Document) error {
	filePath := fs.getFilePath(doc.ID, doc.Name)

	// Проверяем кэш
	if cached, ok := fs.cache.Get(doc.ID); ok {
		w.Header().Set("Content-Type", cached.mime)
		w.Header().Set("ETag", cached.etag)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", cached.size))
		w.WriteHeader(http.StatusOK)
		return nil
	}

	// Получаем метаданные с диска
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return os.ErrNotExist
		}
		return err
	}

	mimeType := doc.Mime
	if mimeType == "" {
		mimeType = mime.TypeByExtension(filepath.Ext(doc.Name))
	}
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	w.WriteHeader(http.StatusOK)
	return nil
}

func (fs *FileStorage) serveFileFromDisk(w http.ResponseWriter, r *http.Request, doc *entity.Document) error {
	lock := fs.getFileLock(doc.ID)
	lock.Lock()
	defer lock.Unlock()

	filePath := fs.getFilePath(doc.ID, doc.Name)

	logrus.Debug("Trying to get file: %v", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return os.ErrNotExist
		}
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Определяем MIME-тип
	mimeType := doc.Mime
	if mimeType == "" {
		mimeType = mime.TypeByExtension(filepath.Ext(doc.Name))
	}
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Кэшируем файл если он небольшой
	if info.Size() < 2*1024*1024 { // 2MB
		data, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		hasher := sha256.New()
		hasher.Write(data)
		etag := hex.EncodeToString(hasher.Sum(nil))

		fs.cache.memoryCache.Store(doc.ID, &CachedFile{
			data:    data,
			size:    info.Size(),
			mime:    mimeType,
			etag:    etag,
			created: time.Now(),
		})
	}

	// Используем эффективную отдачу файлов
	http.ServeContent(w, r, doc.Name, info.ModTime(), file)
	return nil
}

func (c *Cache) Get(id uuid.UUID) (*CachedFile, bool) {
	return c.memoryCache.Get(id)
}

func (mc *MemoryCache) Get(id uuid.UUID) (*CachedFile, bool) {
	mc.RLock()
	defer mc.RUnlock()

	if file, ok := mc.files[id]; ok {
		// Проверяем TTL
		if time.Since(file.created) < cacheTTL {
			return file, true
		}
		// Удаляем просроченный элемент
		go mc.Delete(id)
	}
	return nil, false
}

func (mc *MemoryCache) Store(id uuid.UUID, file *CachedFile) {
	mc.Lock()
	defer mc.Unlock()

	// Если файл слишком большой, не кэшируем
	if int64(len(file.data)) > mc.maxSize {
		return
	}

	// Очищаем место при необходимости
	for len(mc.files) >= mc.maxEntries || mc.totalSize+int64(len(file.data)) > mc.maxSize {
		var oldestKey uuid.UUID
		var oldestTime time.Time
		for key, f := range mc.files {
			if oldestTime.IsZero() || f.created.Before(oldestTime) {
				oldestKey = key
				oldestTime = f.created
			}
		}
		if !oldestTime.IsZero() {
			mc.deleteLocked(oldestKey)
		}
	}

	// Сохраняем файл
	if file.data != nil {
		mc.totalSize += int64(len(file.data))
	}
	mc.files[id] = file
}

func (mc *MemoryCache) Delete(id uuid.UUID) {
	mc.Lock()
	defer mc.Unlock()
	mc.deleteLocked(id)
}

func (mc *MemoryCache) deleteLocked(id uuid.UUID) {
	if file, exists := mc.files[id]; exists {
		if file.data != nil {
			mc.totalSize -= int64(len(file.data))
		}
		delete(mc.files, id)
	}
}
