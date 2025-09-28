# email-app
send email via apixlab contact form 

# Installation go 
download go 
    - wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
extract to /usr/local
    - sudo rm -rf /usr/local/go
    - sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
edit/reload shell profile 
    - vim ~/.bashrc
    - export PATH=$PATH:/usr/local/go/bin
    - source ~/.bashrc
verify go 
    - go version

# Run email-app
build binary
    - go build -o email-app ./cmd/server
run forever with systemd
    - sudo vim /etc/systemd/system/email-app.service
    
    ```
    [Unit]
    Description=Go Email App
    After=network.target

    [Service]
    ExecStart=/opt/apixlab/email-app/email-app
    WorkingDirectory=/opt/apixlab/email-app
    Restart=always
    RestartSec=5
    EnvironmentFile=/opt/apixlab/email-app/configs/.env
    User=www-data
    Group=www-data

    [Install]
    WantedBy=multi-user.target
    ```
save and reload systmd
    - sudo systemctl daemon-reload
    - sudo systemctl start email-app
    - sudo systemctl enable email-app

monitor logs
    - systemctl status email-app
    - sudo journalctl -u email-app -f

restart if needed
    - sudo systemctl restart email-app

stop if need 
    - sudo systemctl stop email-app