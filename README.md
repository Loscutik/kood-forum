# Project name: FORUM

This is a project of a web forum, which contains such things as:

* beautiful and user-friendly interface
* creating user accounts
* password hashing
* usage of cookies and user sessions
* ...and other cool things

Before we start, here you can find the [audit](https://github.com/01-edu/public/tree/master/subjects/forum/audit) questions. Are you ready? Let's go!

## How to run the app?

The project requires Docker, so if you have no Docker on your computer, you need to install it. To do it, please follow this [>>> LINK <<<](https://docs.docker.com/get-docker/). You may also need to install the Docker extension to your VS Code if you use it. Please hold the Docker app opened while you launch the forum project.

The next step is very easy - please open the "forum" folder using your terminal, and run the file named "docker.sh". On a Windows PC, you can do it by using the command `bash docker.sh`, and you can use a command `./docker.sh` for Linux. This will launch a script which will do all the rest needed things to unpack and open the forum.

If, for some case, you want to open the forum project without a Docker, please open the "forum" folder using your terminal, and use a command `go run ./app` to launch the server on port 8080. This command launches a web server written in GoLang. After doing that, if you use Windows PC, you may see an alert, where you need to allow access for Windows Defender Firewall. Please do it.
If you want to launch the app using another port, please just use the "port" flag `go run ./app --port=8080` and replace "8080" with your desired port number.
Now you can go to your web browser, and type the following address in the address bar to open the forum main page: [http://localhost:8080/](http://localhost:8080/) (or just use this link 😉).

## How to use the forum?

So, you are on the main forum page. As a guest, you can see a media feed with the last topics and click on the link "Comments" to open a post page and read comments. The number in brackets next to "Comments" displays how many comments there are in this topic. For example, `Comments (2)` means that there are already 2 comments. You can also see the number of likes and dislikes, and filter posts by categories, but if you want to create a post, write a comment or put your like or dislike, you need to create an account or log in.

To create an account, please click on the "Sign up" button on the right side of the top menu panel. Then, please add your email (it will be visible to other users), choose your username (there will be no option to change it later), choose and confirm your password. If you are not sure if your password is typed correctly, you can use the option "Show password" to see the password which was written in both password fields. When you are ready, please click "Sign up" to create your account. You will be redirected to the page with a button "Sign in". Now you can use this button and log in to your account.

So, you have successfully created an account and logged in. Now you can see your username at the top panel of the website. If you click on the username, you will see an option to see your profile, go to "Settings" to change your email or password, or log out. Don't hesitate to try these options.
Please note, that you will be logged out automatically after 20 minutes of session.

If you want to return to the main page with forum feed, please just click on the website logo.

On the main page, now you can read the posts, put likes or dislikes, write comments. Please note, that if you click on the button near the right corner of the first post, you can open filters. As a registered user, you can not only filter posts by categories, but also filter them by posts created by you, liked by you and disliked by you. Of course, you can mix these filters.

In order to create a new post, please use the button "Create post" at the top panel of the website. Please write a post header, your text, and choose a category of your post below. All the parts are required, and you will be unable to post your topic without filling a post header, the main text and choosing a post category.

The post is posted? You are nice! Now go and write some comments!😉

## The features

#### User-friendly interface.

_Please note, that the interface is not responsive, it's highly recommended to use it on a laptop/PC screen._

We added a lot of features to help user feeling comfortable with our product. Just take a look at some of them:

* The forum follows the principles of UX/UI. The website design is intuitive and meets the expectations. The elements and menus are there where they are intuitively expected to be.
* The error messages are made in red colors, and the notifications of successful actions are made in green color. 
* The "Sign in" menu offers to sign up in case you have not done it, and you have an option to go to the logging menu from the "Sign up" menu.
* The top panel has an option to go to the main page by clicking on a site logo. It disappears when you scroll down and appears again if you touch the top of the page with your cursor.
* Beautiful boxes and toggle buttons instead of ordinary checkboxes.

#### User accounts.

Each user account contains important information:

* Personal username of each user
* Email, which can be found in user's personal profile (you can open it by clicking on his username)
* Account creation date, which can be also found in the profile
* Hashed password. All passwords are stored in a database only encrypted, and it is an additional security feature
* The option to sort posts by author, likes or dislikes (this option is available only for the current user)

#### Other features.

* The website uses cookies to determine the user session. The length of each session is 20 minutes, and you need to log in again after this time.




Created by: **Olena Budarahina** (Gitea username: obudarah), **Kristina Volkova** (Gitea username: Mustkass).