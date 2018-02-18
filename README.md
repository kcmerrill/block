![Block](block.png "Block")

# Block

A really simplistic approach to a CLI launcher. Using a very *cough* sophisticated approach to determining which file/folder/script to run. Many uses, but here are a few:

* Open applications via the cli
* Change directories quickly
* Open folders/code with your favorite editor quickly
* Pause/Play spotify via the cli with just a few characters

# Usage

In order to get the full functionality of `block`, you'll need to setup a simple block function in your bash profile. Before doing so, it's recommended that you try out block a few tries to make sure it fits your needs, and to ensure that the logic used to score the inventory is to your liking. 

```sh
echo 'b () { $(block "$@" | tail -1) ; }' >> ~/.bash_profile
```

### Scripts

Any matches that are found within your `$HOME/block/` directory will be run with `bash -c $HOME/lock/matched.file`. This is incredibly useful for applescripts to pause/play spotify for example(see the `scripts` folder), compose email messages, or any other scripted tasks you can think of. I've included my applescript file. When placed in `$HOME/block/spotify` for me, I can simply run `$ b spot` and that will pause/play my spotify when I'm on the CLI which I am most of the day without needing to do other key commands(my work keyboard doesn't have pause/play buttons. :shame:)

```sh
// music playing
$ b spot
// music pauses
```

### Defaults

By default(outside of the scripts directory outlined above) block will attempt to `open` or `cd` depending on the result set it matched on. So if it scored highest a directory, block will change to that directory, or if it's a file type, it will attempt to `open` it. You can however influence it's decesion by passing in the first param of the search to be the application you want to use. A few examples:

```sh
// If you know you are looking for a super cool directory
$ b cd somesupercoolfuzzymatcheddirectory
// switched to that directory if succesful

// Lets open a known file in vim for easy editing
$ b vim inventory.go
// matches on inventory.go and opens it in vim

// unsure? lazy? block will do what it thinks is right
$ b inventory.go // will `open path/to/inventory.go`
```

## Downloads
![Mac OSX](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/apple_logo.png "Mac OSX") [386](http://go-dist.kcmerrill.com/kcmerrill/block/mac/386) | [amd64](http://go-dist.kcmerrill.com/kcmerrill/block/mac/amd64)

![Linux](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/linux_logo.png "Linux") [386](http://go-dist.kcmerrill.com/kcmerrill/block/linux/386) | [amd64](http://go-dist.kcmerrill.com/kcmerrill/block/linux/amd64)

## Thoughts

This is just a simple POC. So the code is gnarly. Like, really gnarly. No tests, no grouping, it doesn't really make sense. I am however curious to see if I will use this as much as I told myself I will. If I do end up using this as part of my daily workflow, then I want to add a bunch of features. 

* Plugins(like mac's quicksilver for email/music/etc)
* Make it faster
* Refactor
* Add a bunch of tests
* Would be sweet if block could "learn" your behaviors
* I feel like there is some machine learning I could utilize in here. That'd be awesome!
* Attach it to global keystrokes(or wrap your terminal session)