# Golerts
A simple RESTful Notification service written in Go. This will primarily be used with iOS push notifications

## Setup
This application is configured to run on Heroku. As such, there are some steps
to follow to get it running in your local.

1. Remove the `.git` directory and create a fresh one for your project

    ```sh
    $ rm -rf .git
    $ git init .
    ```
2. Install Heroku Toolbelt:

    Via Homebrew:

    ```sh
    $ brew install heroku
    ```

    Or Via [Heroku][1]

3. Create your .env file

    ```sh
    $ cp env.example .env
    ```

4. Update any values in `.env` as needed

5. Run the application locally

    ```sh
    $ go install && heroku local
    ```

## Deploying
Assuming you have the proper heroku app [setup + git remote][2]

1. Remove `vendor` and `Godeps` from .gitignore

2. Save and commit dependencies

    ```sh
    $ godep save
    $ git commit -am "initial import of dependencies for heroku"
    ```

3. Deploy

    ```sh
    $ git push heroku [BRANCH]:master
    ```

## Migrations
You can run the migrations an seeds located in the `migrates` directory to get your database in the most recent state.

NB: After release, new migrate and seeds files should be created to incremental updates

[1]: https://toolbelt.heroku.com/
[2]: https://devcenter.heroku.com/articles/creating-apps