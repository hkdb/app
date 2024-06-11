### BACKGROUND (Probably TL;DR for most)

`If you found this repo via Medium, you have read this already.`

So I personally have a considerable amount of machines. My team and I also deal with a lot of machines at work. When I get a new machine or want to repurpose an existing machine, I prefer to do fresh installs rather than imaging hard drives because among other reasons, imaging carries all the junk/cache over beyond just configs/preferences and it's not only a pain to maintain but also undesirable to keep a big number of images around on storage. Furthermore, handling image change is very time consuming. Plus, I mean, there's something to be said about a clean install if you know what I mean. 

Sure you can carry over just the home directory minus your .cache directory which mostly takes care of the junk and cache but how will you get all the same software installed that your home directory references? And yes, you can see what you manually installed with each package manager on your older computer running Debian based distros if you just do a little googling and then ninja your way through the output to derive long lines of packages to install manually on the new machine but God forbid you are on a rpm based system… I am not even aware of a way to see what was manually installed. I guess back to googling? What about Macs and Windows which I often use Homebrew and Scoop on?

Sometimes, I also need to use a different distro and more often than not, the distro at hand is still based on 1 of the 3 grandfathers (Debian, Redhat, Arch). It's also not necessarily always feasible to maintain an ansible playbook since often times, people forget they even installed something new while they were in the heat of their workflow. So on your next computer, you go through that whole troubleshooting journey again just to end up fixing the same problem by installing the same missing library manually.

How about package management standardization? Unless you are a purist which I am not dismissing the fact that purists do exist, chances are, you use the native package manager where you can and Flatpak just because your favorite software's developer doesn't want to deal with cross-distro hell and maybe even *evil music*… snap for certain server side packages that work better or are just easier to setup with snap.

Last but not least, it's only natural that I get slowed down during work because I accidentally typed `sudo apt install` on an Arch based distro and then having to Ctrl-C half way through it to finally type `sudo pacman -S` just so something trivial is installed. I wonder how many hours this type of incident collectively took in my 20+ years of IT. Sure it's just a second or 2 MAX each time but I am willing to bet that if I counted them every time, the number would end up significant.

![wrong](assets/wrong.png)

And so......... out came this project that was born over some time off from work during the winter holidays in 2023. It aims to solve the problem by:

- Providing the same command structure across OS/distros and package managers to reduce the brain teasing task of remembering which command structure is for which package manager
- Automatically remember what repos you added and what packages you installed on each package manager
- Provide a single command upgrade of all the software across all package managers used
- Provide a way to allow for a single command install of everything you have ever installed on another computer if it's the same distro base or OS
