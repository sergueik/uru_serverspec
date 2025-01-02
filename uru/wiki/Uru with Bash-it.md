[Bash it](https://github.com/revans/bash-it) is a project to make customizing your Bash terminal an easy process. One of the nice features is the ability to see which Ruby version you're currently using. Unfortunately, this doesn't work out of the box with Uru. It is possible, however, to use the included [chruby](https://github.com/postmodern/chruby) plugin to get a nice looking prompt.

Here's my prompt without the plugin:

```
 Macbook in ~/Play/uru
± |master ✗| → 
```

## Steps

1. Follow the instructions to install *Bash it* https://github.com/revans/bash-it#install
2. Install *Chruby* but don't update the `.bash_profile` or `.bashrc` files https://github.com/postmodern/chruby#install
3. Enable the "chruby" plugin for *Bash it* with `bash-it enable plugin chruby`
4. Install *Uru* and make sure that it's loaded in the `.bash_profile` or `.bashrc`
5. Open a new terminal session. Just sourcing may not reload *Bash it* properly

If all goes well, your terminal should look like the following (I'm using the "envy" *Bash it* theme):

```
rubinius 2.1.1 Macbook in ~/Play/uru
± |master ✗| → 
```