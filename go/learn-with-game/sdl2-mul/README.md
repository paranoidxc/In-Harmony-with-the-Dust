// install sdl2.dmg
https://www.libsdl.org/
sudo cp -r /Volumes/SDL2/SDL2.framework /Library/Frameworks

brew install pkg-config

```
brew install sdl2{,_image,_mixer,_ttf,_gfx} pkg-config
```

brew install sdl2_ttf pkg-config




go get github.com/veandco/go-sdl2/sdl

go get -v github.com/veandco/go-sdl2/sdl
go get -v github.com/veandco/go-sdl2/img
go get -v github.com/veandco/go-sdl2/mix
go get -v github.com/veandco/go-sdl2/ttf
go get -v github.com/veandco/go-sdl2/gfx

go run sdl2.go
