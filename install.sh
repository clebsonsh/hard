#!/usr/bin/env bash

echo "Installing Hard..."

# check if git is installed
if ! command -v git >/dev/null; then
  echo "Git is not installed. Please, install git and try again."
  exit 1
fi

# check if docker is installed
if ! command -v docker >/dev/null; then
  echo "Docker is not installed. Please, install docker and try again."
  exit 1
fi

# check if docker-compose or docker compose is installed
if ! command -v docker-compose >/dev/null; then
  if ! command -v docker compose >/dev/null; then
    echo "Docker Compose is not installed. Please, install docker-compose and try again."
    exit 1
  fi
fi

# set HARD_PATH
HARD_PATH="/home/${USER}/.hard"

# check if HARD_PATH exists
if [ -d $HARD_PATH ]; then
  # pull hard repository
  cd $HARD_PATH; git pull
else
  # clone hard repository
  git clone https://github.com/clebsonsh/hard.git $HARD_PATH
fi


# copy .env.example to .env if .env
cp --update=none $HARD_PATH/.env.example $HARD_PATH/.env

# replace USER and USER_ID in .env for the current user
sed -i "s/USER=hard/USER=$(echo $USER)/" $HARD_PATH/.env
sed -i "s/USER_ID=1001/USER_ID=$(echo $UID)/" $HARD_PATH/.env

# source .env
. $HARD_PATH/.env

echo "If you want to change the default values, please, edit the .env file in $HARD_PATH."

# check if WWW_PATH exists and create it
if [ ! -d $WWW_PATH ]; then
  mkdir -p $WWW_PATH
fi

# check if .bashrc exists and have ~/.local/bin in the PATH
if [ -f ~/.bashrc ]; then
  if ! grep -q ".local/bin" ~/.bashrc; then
    echo "export PATH=\$PATH:~/.local/bin" >>~/.bashrc
  fi
fi

# check if .zshrc exists and have ~/.local/bin in the PATH
if [ -f ~/.zshrc ]; then
  if ! grep -q ".local/bin" ~/.zshrc; then
    echo "export PATH=\$PATH:~/.local/bin" >>~/.zshrc
  fi
fi

# check if config.fish exists and have ~/.local/bin in the PATH
if [ -f ~/.config/fish/config.fish ]; then
  if ! grep -q ".local/bin" ~/.config/fish/config.fish; then
    echo "set -x PATH \$PATH ~/.local/bin" >>~/.config/fish/config.fish
  fi
fi

# check if ~/.local/bin exists
if [ ! -d ~/.local/bin ]; then
  mkdir -p ~/.local/bin
fi

# check if ~/.local/bin/hard exists
if [ -f ~/.local/bin/hard ]; then
  rm ~/.local/bin/hard
fi

# create a symbolic link to hard binary
ln -s -f -t ~/.local/bin $HARD_PATH/bin/hard

# give permission
chmod +x $HARD_PATH/hard

echo "Hard installed successfully!"
echo "Please, restart your terminal."
