package main

import (
  "context"
  "fmt"
  "fyne.io/fyne"
  "fyne.io/fyne/app"
  "fyne.io/fyne/canvas"
  "fyne.io/fyne/dialog"
  "fyne.io/fyne/layout"
  "fyne.io/fyne/widget"
  "github.com/shibukawa/configdir"
  "image/color"
  "log"
  "time"
)

var oauthStateToken string
var port int

var appId = "org.develar.lr-backup"

func main() {
  configDirs := configdir.New("", appId)
  cache := configDirs.QueryCacheFolder()
  // check, is there saved refresh token to exchange it to an access token
  token := readToken(cache)
  if len(token) == 0 {
    // no saved refresh token - open window and ask to sign in
    openSignInWindow()
  }
}

func auth(a fyne.App) error {
  server, tokenChannel, err := startServer()
  if err != nil {
    return err
  }

  authCodeUrl, err := generateSignInUrl()
  if err != nil {
    return err
  }

  err = a.OpenURL(authCodeUrl)
  if err != nil {
    return err
  }

  token := <-tokenChannel
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  err = server.Shutdown(ctx)
  if err != nil {
    LogError("cannot shutdown server: %v", err)
  }
  log.Printf("got code: %s", token.AccessToken)
  return nil
}

func openSignInWindow() {
  windowWidth := 480
  headerSize := fyne.NewSize(windowWidth, 59)
  headerBackgroundColor := color.RGBA{
    R: 57,
    G: 57,
    B: 57,
    A: 255,
  }
  //contentBackgroundColor := color.RGBA{
  //  R: 66,
  //  G: 66,
  //  B: 66,
  //  A: 255,
  //}

  a := app.NewWithID(appId)
  w := a.NewWindow("Sign In")
  w.Resize(fyne.NewSize(windowWidth, 651))
  w.SetFixedSize(true)

  headerBackground := canvas.NewRectangle(headerBackgroundColor)
  headerBackground.Resize(headerSize)
  headerBackground.SetMinSize(headerSize)
  headerText := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), layout.NewSpacer(), canvas.NewText("Lightroom Backup", color.White), layout.NewSpacer())
  headerText.Resize(headerSize)
  //header := fyne.NewContainer(headerBackground, headerText)

  //left := layout.NewSpacer()
  //middle := canvas.NewText("content", color.White)

  //text4 := canvas.NewText("centered", color.White)
  //centered := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), text4, layout.NewSpacer())
  //	myWindow.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), container, centered))

  hello := widget.NewLabel("This app will help you backup all of your edits Lightroom photos.")
  //content := fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil), header, hello)

  progressBar := widget.NewProgressBarInfinite()
  progressBar.Hide()
  var centeredButton *fyne.Container
  var button *widget.Button
  button = widget.NewButton("Sign In", func() {
    button.SetText("Signing in…")
    button.Disable()
    progressBar.Show()
    go func() {
      err := auth(a)
      if err != nil {
        button.Enable()
        button.SetText("Sign In")
        progressBar.Hide()
        dialog.ShowError(fmt.Errorf("internal error: %w", err), w)
      }
    }()
  })
  centeredButton = fyne.NewContainerWithLayout(layout.NewHBoxLayout(), layout.NewSpacer(), button, layout.NewSpacer())

  content := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), layout.NewSpacer(), hello, centeredButton, progressBar, layout.NewSpacer())

  w.SetContent(content)

  w.ShowAndRun()
}
