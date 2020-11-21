# **Market-Place**

## **overview**:
mini market place integrate with raja ongkir api

## **Demo Link**:
http://hokdre.com

## **User**:
1. **Customer**: <br/>
   email : hadinw@gmail.com <br/>
   password : password!23Z <br/>
2. **Merchant**:<br/>
   1. Merchant Mike (Laptop)<br/>
      email : mukminmike@gmail.com <br/>
      password : password!23Z <br/>
   2. Merchant Jainal (Handphone)<br/>
      email : jainal@gmail.com<br/>
      password : password!23Z <br/> 
3. **Admin**:<br/>
   email : superadmin@gmail.com <br/>
   password : password!23Z <br/>

## **Teknologi Yang Digunakan** :
1. ### **Server:**
   gorilla/mux <br/>
   Gorrila mux pada project ini, saya gunakan untuk membuat http server.
2. ### **Authentication:**
   JWT <br/>
   JWT pada project ini, saya gunakan sebagai encoder dan decoder token untuk proses autentikasi dan otorisasi.
3. ### **Cache:**
   Redis <br/>
   Redis pada project ini, saya gunakan untuk melakukan cache ongkos kirim hasil fetch API Raja Ongkir dan fitur autocomplete input city.
4. ### **Database**
   1. MongoDB <br/>
      Mongodb pada project ini, saya gunakan sebagai database utama. 
   2. ElasticSearch <br/>
      Elasticsearch pada project ini, saya gunakan untuk melakukan pencarian data produk dan merchant, selain data produk dan merchant tidak ada data lain yang disimpan.
5. ### **Deployment**
   Google VM dan storage bucket
6. ### **APIs**
   1. Raja Ongkir <br/>
      Api raja ongkir, saya gunakan untuk melakukan fetch data ongkos kirim serta ketersedian layanan pengiriman yang ada dari kota satu ke kota yang lain. (Pengiriman POS : yang dicover oleh raja ongkir)
7. ### **Other :**
   1. Kafka <br/>
      Kafka pada project ini, saya gunakan sebagai perantara syncronisasi data mongodb dengan elasticsearch.  dimana setiap event yang terjadi pada database mongodb akan disubscribe melalui kafka connect dengan bantuan plugin mongo kafka connect, lalu event akan disimpan pada topic di kafka server.
      Pemilihan penggunaan kafka sendiri, karena ketika melakukan pencarian plugin mongodb source di logstash saya tidak menemukan plugin offical ataupun plugin pihak ketiga yang masih di maintenance. 
   2. Logstash <br/>
      Logstash pada project ini, saya gunakan sebagai mensubsribe data event dari kafka topic di server, lalu data event akan ditransform dan kemudian disimpan pada elasticsearch. 


## **Development Problem**:
   Kesalahan yang saya buat selama melakukan development projek ini ialah diawal projek tidak melakukan riset yang cukup terkait API pihak ke 3 yang digunakan. Hal ini menyebabkan desain database yang saya buat tidak sesuai dengan data yang didapat dari API pihak ke 3 dan data yang dibutuhkan oleh pihak API ke 3 juga tidak tersedia pada database yang saya buat.  Ini merupakan hal yang sangatlah fatal, database yang telah dibuat harus diubah kembali, sehingga semua fitur yang telah dikerjakan juga harus di sesuaikan ulang. 

## **Requirement Sebelum Menjalankan Program** :
   1. gcp.json 
      file json sebagai admin storage bucket google cloud platform, letakan file ini di root projek
   2. RAJA_ONGKIR_KEY (string)
      set pada environment variable 
   3. JWT_SECRET (string)
      set pada environment variable  
   4. MONGO_URL (host:port)
      set pada environment variable  
   5. ELASTIC_URL (host:port)
      set pada environment variable  
   6. REDIS_URL (host:port)
      set pada environment variable  

## **cara menjalankan program dengan instalasi biasa** :
   1. jalankan aplikasi pendukung mongodb, redis, elasticsearch, zookeeper,dan kafka server seperti biasa.
   2. Copy monggo connector plugin pada /deployment_setting/manual/plugins/kafka pada project ini ke /usr/local/share/kafka/plugins
   3. jalankan kafka connect stand alone dengan configuration file yang ada pada /deployment_setting/connect-standalone.properties
      ```
      ${instalasi_path}/bin/connect-standalone.sh  ${instalasi_path}/config/connect-standalone.properties /usr/local/share/kafka/plugins/etc/MongoSourceConnector.properties
      ```
   4. jalankan logstash dengan configuration file yang ada pada /deployment_setting/logstash.conf
      ```
      ${instalasi_path}/bin/logstash -f ${dir}/market-place/deployment_settings/manual/logstash.conf
      ``` 

## **cara menjalankan program dengan instalasi docker : FAILED** :

| WARNING:  still FAILED :'( , menjalankan kafka connect stand alone didocker sangat laggy dan sering kali process killed, saat ini saya menjalankan kafka connect stand alone menggunakan satu image yang sama dengan kafka server yaitu melalui docker exec. Alternatif kedepan mungkin dengan memisahkan antara image kafka server dan kafka connect stand alone|
| --- |

1. Menjalankan Program
   ```
   docker-compose up -d
   ```
2. Menjalankan kafka connect
   ```
   //1. masuk kedalam bash container kafka
   docker exec -it kafka bash
   //2. start manual kafka connect
   /opt/kafka/bin/connect-standalone.sh /opt/kafka/config/connect-standalone.properties usr/local/share/kafka/plugins/mongodb-kafka-connect-mongodb-1.2.0/etc/MongoSourceConnector.properties
   ```
3. Monitoring logs
   ```
   docker-compose logs -f web
   docker-compose logs -f mongo
   docker-compose logs -f redis
   docker-compose logs -f zookeeper
   docker-compose logs -f kafka
   docker-compose logs -f logstash
   docker-compose logs -f elastic
   ```

## **Migrations REST API**:
```
curl -X GET 'http://localhost:80/elastic-product-index'
curl -X GET 'http://localhost:80/elastic-merchant-index'
```

## **Seeder REST API**:
| WARNING:  urutan seeder harus sesuai|
| --- |
```
curl -X GET 'http://localhost:80/seed-admin'
curl -X GET 'http://localhost:80/seed-shipping'
curl -X GET 'http://localhost:80/seed-customer'
curl -X GET 'http://localhost:80/seed-merchant'
curl -X GET 'http://localhost:80/seed-product'
```
