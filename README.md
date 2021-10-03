# Library Backend

### Create Service File
1. To run the library at startup create a service file called libary.service in your current directory.
Copy and paste this file template, and make sure to verify your ExexStart and WorkingDirectory are correct.
```
[Unit]
Description=My service
After=network.target

[Service]
ExecStart=/usr/local/go/bin/go run main.go &
WorkingDirectory=/home/pi/library_backend
StandardOutput=inherit
StandardError=inherit
Restart=always
User=pi

[Install]
WantedBy=multi-user.target
```
2. Copy the service file to lib/systemd
```
sudo cp libary.service /lib/systemd/system/library.service
```


To get logs of the running proccess
```
sudo journalctl -fu library.service
```

To change databases in postgres console:
```bash
\c library
```