## APP CONFIG DIRECTORY STRUCTURE

The config directory app creates on your machine has the following structure to "remember what you installed and removed". Some directories may not exist for you since they are not created until you install something with that specific package manager.

```
~/.config/app
      |_ settings.conf
      |_ packages
          |
          |_ appimage.json
          |
          |_ apt.json
          |
          |_ deb.json - a list of packages
          |
          |_ brew.json
          |
          |_ dnf.json
          |
          |_ pkg.json
          |
          |_ zypper.json
          |
          |_ flatpak.json
          |
          |_ go.json
          |
          |_ pip.json
          |
          |_ cargo.json
          |
          |_ pacman.json
          |
          |_ local
          |   |_ deb
          |   |   |_ <package name>.json
          |   |   |_ <package file name>.deb
          |   |   
          |   |_ rpm
          |   |   |_ <package name>.json
          |   |   |_ <package file name>.rpm
          |   |
          |   |_ appimage
          |       |_ <package name>.json>
          |       |_ <package name>
          |           |_ <package name>.AppImage
          |           |_ <package>.desktop
          |           |_ <package>.png
          |
          |_ repo
          |   |_ apt.json
          |   |_ dnf.json
          |   |_ pacman.json
          |   |_ zypper.json
          |   |_ yay.json
          |   |_ flatpak.json
          |   |_ cargo.json
          |   |_ Channel
          |   |   |_snap
          |   |      |_ <package name>.json
          |   |
          |   |_ local
          |       |_ apt
          |       |   |_ <repo name>.sh
          |       |
          |       |_ dnf
          |       |   |_ <repo name>.sh/json
          |       |       
          |       |_ pacman
          |       |   |_ <repo name>.sh
          |       |
          |       |_ zypper
          |       |   |_ <repo name>.sh
          |       |
          |       |_ flatpak
          |       |   |_ <repo name>.json
          |       |
          |       |_ snap
          |       |   |_ <package name>.json 
          |       |
          |       |_ cargo
          |           |_ <package name>.json 
          |
          |
          |_ rpm.json
          |
          |_ snap.json
          |
          |_ yay.json
          |
          |_ paru.json
```

You can of course choose to work backwards and mannually edit/compose the data inside this directory but you better know what you are doing or something really wrong could happen. Be warned.... 

