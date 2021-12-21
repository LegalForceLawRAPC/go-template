# go-template

Steps to Deploy:
```
sudo docker build -t legalforce/go-template .
sudo docker run -d -p 3000:3000 --name template --env-file .env legalforce/go-template
```

Steps to Redeploy 
```
sudo docker build -t legalforce/go-template .
sudo docker stop template
sudo docker rm template
sudo docker run -d -p 3000:3000 --name template --env-file .env legalforce/go-template
```
