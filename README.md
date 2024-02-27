# Forum Modern

In forum image upload, registered users have the possibility to create a post containing an image as well as text.

When viewing the post, users and guests should see the image associated to it.
There are several extensions for images like: JPEG, SVG, PNG, GIF, etc. In this project you have to handle at least JPEG, PNG and GIF types.

The max size of the images to load should be 20 mb. If there is an attempt to load an image greater than 20mb, an error message should inform the user that the image is too big.

## Overview
This project is a forum application designed to facilitate communication and interaction between users. It includes features such as authentication with tokens, communication between users, categorizing posts, liking and disliking posts and comments, as well as filtering posts. Users can also edit and delete their own posts and comments.

## Features
- **Authentication with Tokens**: Users can authenticate themselves using tokens for secure access to the forum.
- **Communication Between Users**: Users can interact with each other through posts and comments.
- **Associating Categories to Posts**: Posts can be categorized to organize discussions effectively.
- **Liking and Disliking Posts and Comments**: Users can express their opinions on posts and comments by liking or disliking them.
- **Filtering Posts**: Users can filter posts based on categories or other criteria to find relevant content.
- **Editing and Deleting Posts and Comments**: Users have the ability to edit and delete their own posts and comments.

## Technologies Used
- **Backend**: Golang, sql
- **Database**: SQLite
- **Authentication**: [UUID](https://github.com/gofrs/uuid), [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- **Frontend**: HTML, CSS, JavaScript

## Installation
1. Clone the repository: `git clone <repository_url>`
2. Navigate to the project directory: `cd forum-project`
3. Run Locally with makefile: ```make run```
4. And go to the web page: `http://localhost:8080`
5. Run docker: 
```
docker build -t forum .
docker run -it -p 8080:8080 forum
make dStop
make dDelete
make dDeleteImages
```

## Usage
1. Register a new account or login with existing credentials.
2. Explore categories or create new posts within categories.
3. Interact with other users by commenting on their posts or liking/disliking their content.
4. Utilize the filtering options to find specific posts or categories.
5. Manage your own posts and comments by editing or deleting them as needed.


## Contact
For any inquiries or support regarding this project, feel free to contact [Elnur Bauyrzhan tgm: @EL_n_UR 
Discord: https://discordapp.com/users/762269394034229248/ 
gitea: https://01.alem.school/git/ebauyrzh].
