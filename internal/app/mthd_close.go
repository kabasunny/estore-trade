// internal/app/app.go
package app

// Close はアプリケーションの終了処理を行う
func (app *App) Close() {
	if app.DB != nil {
		app.DB.Close()
	}
	if app.Logger != nil {
		app.Logger.Sync()
	}
}
