# Bring down nickanderssohn.com before running this
docker run -p 80:80 -p 443:443 -it --rm --name certbot -v "/etc/letsencrypt:/etc/letsencrypt" -v "/var/lib/letsencrypt:/var/lib/letsencrypt" certbot/certbot --standalone renew
