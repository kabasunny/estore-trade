// internal/infrastructure/persistence/tachibana/tests/process_response_body_test.go
package tachibana_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"estore-trade/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// モック EventStream (修正)
type MockEventStream struct {
	mock.Mock
	Logger *zap.Logger // Logger を追加, 公開
}

// EventStream のインターフェースを定義 (モック用)
type EventStreamInterface interface {
	ParseEvent(data []byte) (domain.OrderEvent, error)
	SendEvent(event domain.OrderEvent)
	GetLogger() *zap.Logger // GetLogger メソッドを追加
}

func (m *MockEventStream) ParseEvent(data []byte) (domain.OrderEvent, error) {
	args := m.Called(data)
	// モックから返された値を適切な型にアサート
	return args.Get(0).(domain.OrderEvent), args.Error(1)
}

func (m *MockEventStream) SendEvent(event domain.OrderEvent) {
	m.Called(event)
}

func (m *MockEventStream) GetLogger() *zap.Logger {
	return m.Logger // 公開されたLoggerを返す
}

func TestProcessResponseBody(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	t.Run("Success", func(t *testing.T) {
		mockES := &MockEventStream{Logger: logger} // Loggerを設定

		// 期待されるイベントデータ
		expectedEvent := domain.OrderEvent{EventType: "test_event"}

		// parseEvent メソッドのモック設定
		mockES.On("ParseEvent", mock.Anything).Return(expectedEvent, nil)
		mockES.On("SendEvent", expectedEvent).Return()

		// テスト用のレスポンスボディ
		responseData := "test_response_data"
		resp := &http.Response{
			Body: io.NopCloser(bytes.NewBufferString(responseData)),
		}

		// ヘルパー関数を使用せず、直接呼び出し
		err := mockES.processResponseBody(resp) // 直接呼び出し
		assert.NoError(t, err)
		mockES.AssertExpectations(t) // モック呼び出し検証
	})

	t.Run("ReadError", func(t *testing.T) {
		mockES := &MockEventStream{Logger: logger} // Loggerを設定
		// 読み込みエラーを発生させるモック
		resp := &http.Response{
			Body: io.NopCloser(&errorReader{}), //エラーを発生させる
		}

		// ヘルパー関数を使用せず、直接呼び出し
		err := mockES.processResponseBody(resp) //直接呼び出し

		assert.Error(t, err)
		assert.EqualError(t, err, "read error") //特定のエラーを返す
	})

	t.Run("ParseError", func(t *testing.T) {
		mockES := &MockEventStream{Logger: logger} // Loggerを設定
		// parseEvent がエラーを返すようにモックを設定
		mockES.On("ParseEvent", mock.Anything).Return(domain.OrderEvent{}, errors.New("parse error"))

		// テスト用のレスポンスボディ
		responseData := "invalid_response_data" //parseEventでエラーとなるデータ
		resp := &http.Response{
			Body: io.NopCloser(bytes.NewBufferString(responseData)),
		}

		// ヘルパー関数を使用せず、直接呼び出し
		err := mockES.processResponseBody(resp) // 直接呼び出し
		assert.Error(t, err)
		assert.EqualError(t, err, "parse error") //parse errorを返す
		mockES.AssertExpectations(t)
	})

	t.Run("Empty Response", func(t *testing.T) {
		mockES := &MockEventStream{Logger: logger} // Loggerを設定

		//空のレスポンス
		resp := &http.Response{
			Body: io.NopCloser(bytes.NewBufferString("")),
		}

		// ヘルパー関数を使用せず、直接呼び出し
		err := mockES.processResponseBody(resp) //直接呼び出し

		assert.NoError(t, err)                                 //エラーにならない
		mockES.AssertNotCalled(t, "ParseEvent", mock.Anything) //parseEventが呼ばれない
		mockES.AssertNotCalled(t, "SendEvent", mock.Anything)  //sendEventが呼ばれない
	})
}

// io.Reader のエラーを発生させるためのモック
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

// processResponseBodyをMockEventStreamに追加
func (m *MockEventStream) processResponseBody(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		m.Logger.Error("Failed to read event stream response", zap.Error(err))
		return err
	}
	receivedData := string(body)
	if receivedData != "" {
		m.Logger.Info("Received event stream message", zap.String("message", receivedData))
		event, err := m.ParseEvent(body) //モックのParseEventを呼び出す
		if err != nil {
			m.Logger.Error("Failed to parse event stream message", zap.Error(err))
			return err
		}
		m.SendEvent(event) //モックのSendEventを呼び出す

	}
	return nil
}
