package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestConnectWithMockDB(t *testing.T) {
	// Mock veritabanı bağlantısı oluştur
	mockDBConnection := new(MockDBConnection)

	// Mock bir gorm.DB nesnesi oluştur
	expectedDB := &gorm.DB{}

	// Expect "DBConnection" işlevinin çağrılmasını ve mock DB nesnesi döndürmesini bekliyoruz
	mockDBConnection.On("DBConnection").Return(expectedDB)

	// Connect işlevini çağır
	ConnectWithDBConnection(mockDBConnection)

	// DB değişkeninin beklenen değeriyle aynı olup olmadığını kontrol et
	assert.Equal(t, expectedDB, DB)

	// Mock beklentilerini kontrol et
	mockDBConnection.AssertExpectations(t)

	fmt.Println(DB)
}

// Bu, gerçek DB bağlantısını oluşturan Connect işlevinizdir
func ConnectWithDBConnection(conn DBConnector) {
	DB = conn.DBConnection()
}

// DB bağlantısı için bir arayüz tanımlayın (gerçek ve mock bağlantıları için kullanılabilir)
type DBConnector interface {
	DBConnection() *gorm.DB
}

// Mock veritabanı bağlantısı
type MockDBConnection struct {
	mock.Mock
}

// DBConnection işlevini mock olarak uygula
func (m *MockDBConnection) DBConnection() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}
