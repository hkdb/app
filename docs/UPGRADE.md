## UPDATE/UPGRADE

Run the below command in the terminal:

```
app -m app update
```

You can also update app by compiling it yourself:

Tracking versioned release:

1. cd back into the repo whereever you put it. `~/.config/app` if you took my recommendation 
2. `git pull`
3. `git checkout <version tag>`
3. `./update.sh` # FreeBSD requires bash to be installed

Tracking main branch:

1. cd back into the repo whreever you put it. `~/.config/app` if you took my recommendation 
2. `git pull`
3. `./update.sh` # FreeBSD requires bash to be installed

