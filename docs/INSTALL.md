## INSTALL

Run the below command in the terminal:

```
curl -sL https://hkdb.github.io/app/getapp.sh | bash
```

You can also install app by compiling yourself:

1. Make sure all the package managers you want app to manage and are installed and configured properly
2. Install `git` manually if it's not already installed
3. cd into a directory of choice where you want to keep the app repo. If you are an end user and can't decide, I suggest `~/.config` (`mkdir ~/.config` if it doesn't already exist)
4. git clone https://github.com/hkdb/app.git`
5. `cd app`
6. Optionally `git checkout <version>`
7. `./install.sh` # Note, for FreeBSD, bash must first be installed.

