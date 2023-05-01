 #!/bin/bash
echo -e "Building docker image and container\n\n"
echo "---- To stop server, press CTRL+C----"
echo "------ Press Enter to Continue ------"
echo -e "\n"
read "part"
docker image build -f Dockerfile -t forum-image .
docker container run -p 8080:8080 --detach  --name forum-container forum-image
echo
echo "**----------------------------------------------**"
echo "  Server is running at http://www.localhost:8080/"
echo "**----------------------------------------------**"
echo "Here on can access to files in the container. "
echo "DB is in /forum/forumDB.db"
read -p "press Enter to continue"
echo "**----------------------------------------------**"
echo -e "\nYou are in the /forum directory now"
echo -e "To check DB use the command:\n"
echo "sqlite3 -table forumDB.db"
echo -e "\nTo exit, type:\n"
echo -e "exit\n\n"
echo "===============>"
docker exec -it forum-container sh
echo -e "\n\n\n"
echo -e "Removing Docker image and container from your system\n\n"
docker rm -f forum-container
docker image rm forum-image
echo "**------------------------------------------------------------------**"
echo "DONE"