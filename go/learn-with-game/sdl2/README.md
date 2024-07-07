// install sdl2.dmg
https://www.libsdl.org/
sudo cp -r /Volumes/SDL2/SDL2.framework /Library/Frameworks

brew install pkg-config

go get github.com/veandco/go-sdl2/sdl

go run sdl2.go