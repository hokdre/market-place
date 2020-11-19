package seed

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/logic"
)

const googleStorageAPI string = "https://storage.googleapis.com/ecommerce_s2l_assets/"

func SeedProduct(productUsecase logic.ProductUsecase, merchants []domain.Merchant) ([]domain.Product, error) {
	log.SetOutput(os.Stdout)
	log.Println("Seed Product: starting!")

	var products []domain.Product
	var errs []error

	var wgSeedProduct sync.WaitGroup
	wgSeedProduct.Add(2)
	for _, merchant := range merchants {
		credential := domain.Credential{
			MerchantID: merchant.ID,
		}
		switch merchant.ID {
		case MerchantMikeID:
			go func() {
				defer wgSeedProduct.Done()
				productsMike, err := seedProductInMikeMerchant(productUsecase, credential)
				if err != nil {
					log.Printf("Seed Product: failed, cause %s \n", err)
					errs = append(errs, err)
					return
				}
				products = append(products, productsMike...)
			}()
		case MerchantJainalID:
			go func() {
				defer wgSeedProduct.Done()
				productJainal, err := seedProductInJainalMerchant(productUsecase, credential)
				if err != nil {
					log.Printf("Seed Product: failed, cause %s \n", err)
					errs = append(errs, err)
					return
				}
				products = append(products, productJainal...)
			}()
		}
	}

	wgSeedProduct.Wait()
	if len(errs) != 0 {
		log.Printf("Seed Product: failed, cause %s \n", errs[0])
		return []domain.Product{}, errs[0]
	}

	log.Println("Seed Product: finish!")
	return products, nil
}

func seedProductInMikeMerchant(productUsecase logic.ProductUsecase, credential domain.Credential) ([]domain.Product, error) {
	log.SetOutput(os.Stdout)
	log.Println("Seed Product Merchant Mike: starting!")

	var products []domain.Product
	var errs []error

	var seedProductInMerchantMike sync.WaitGroup
	seedProductInMerchantMike.Add(10)

	//product 1
	go func() {
		defer seedProductInMerchantMike.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		samsungGalaxyFlex := adapter.ProductCreateInput{
			Name: "SAMSUNG GALAXY BOOK FLEX 13.3\" QLED Intel i5 16GB 512GB Win10",
			Category: domain.Category{
				Top:       "komputer & laptop",
				SecondSub: "komputer & laptop",
				ThirdSub:  "laptop",
			},
			Etalase: "laptop",
			Tags: []string{
				"laptop bagus",
				"laptop samsung",
				"laptop ringan",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  30,
			Height: 30,
			Long:   30,
			Description: `Deskripsi SAMSUNG GALAXY BOOK FLEX 13.3 QLED Intel i5 16GB 512GB Win10
			Samsung Galaxy Book Flex NT930QCG-K58 13.3 QLED Intel i5 16GB 512GB

			The Galaxy Book Flex feature Wireless PowerShare. It's a reverse wireless charging feature that Samsung introduced with the Galaxy S10 this year. You can change a Qi-compatible wireless charging device using the laptop’s own battery. To do that, simply place a compatible device on the trackpad.

			The world's first QLED display
			Model Number : NT930QCG-K58
			COLOUR : SESUAI DI PHOTO

			Specifications

			Processor : Intel® Core™ i5-1035G4 Processor
			OS : Windows 10 Home
			Memory : 16 GB LPDDR4x Memory (On BD 16 GB)
			Storage : 512 Gb NVMe SSD / 1 SSD Slot
			Graphic : Intel Iris Plus Graphics
			Display : 33.7 cm (13.3 inch) FHD QLED Wide Viewing Angle Display (1920 x 1080)
			Camera : 720p HD camera
			Network : Bluetooth v5.0 / Wi-Fi 6 (Gig+), 802.11 ax 2x2
			Audio : AKG Stereo speakers (Max 5W x2)
			Weight : 1.16 kg
			Size : 302.6 x 202.9 x 12.9 mm
			Input device : Pebble keyboard (with backlit) / Clickpad, Touch Screen, S Pen
			Security Features : TPM, Fingerprint recognition
			Port : 2 Thunderbolt 3 / 1 USB-C / 1 Headphone output / mic input combo / UFS & MicroSD Combo
			Volage : 100-240 V`,
			Price: 36999000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, samsungGalaxyFlex)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-galaxy-flex-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-galaxy-flex-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-galaxy-flex-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-galaxy-flex-4.png"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-galaxy-flex-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 2
	go func() {
		defer seedProductInMerchantMike.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		samsungNfNc := adapter.ProductCreateInput{
			Name: "Samsung Np Nc 108",
			Category: domain.Category{
				Top:       "komputer & laptop",
				SecondSub: "komputer & laptop",
				ThirdSub:  "laptop",
			},
			Etalase: "laptop",
			Tags: []string{
				"laptop bagus",
				"laptop samsung",
				"laptop ringan",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  30,
			Height: 30,
			Long:   30,
			Description: `Deskripsi Samsung Np Nc 108 Ram ddr3 2gb hdd 320gb laptop leptop minus
			Netbook Samsung
			seri lengkapnya NP-NF208
			Lcd / layar 10,1" bening.
			Intel atom.
			Ram ddr3 2gb. ( lebih cepat dari ddr2)
			hdd /penyimpanan 320gb. Lega.
			Warna putih kombinasi hitam.

			Fitur tambahan: Bluetooth, wifi, usb port. audio.

			bekas.
			Yang didapat: Unit netbooknya, Adaptor original, bonus mouse. lainnya gak ada.

			Packing aman. makai buble warp dan kardos tebal.

			-Minusnya: Batre drop total.harus colok charger.
			Tombol CTRL mati. Fungsi CTRL dipindahkan ke SHIFT sebelah kiri.

			Stok produk tersedia. Silahkan langsung order kak. Terimakasih.
			`,
			Price: 4000000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, samsungNfNc)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-np-nc-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-np-nc-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-np-nc-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-np-nc-4.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-np-nc-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 3
	go func() {
		defer seedProductInMerchantMike.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		samsungNfNc := adapter.ProductCreateInput{
			Name: "Samsung 15.6\" Notebook 5 NT500R5W-XD7S i7-7500U 2.7GHz DDR4 8GB SSD 25",
			Category: domain.Category{
				Top:       "komputer & laptop",
				SecondSub: "komputer & laptop",
				ThirdSub:  "laptop",
			},
			Etalase: "laptop",
			Tags: []string{
				"laptop bagus",
				"laptop samsung",
				"laptop ringan",
				"samsung notebook",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  30,
			Height: 30,
			Long:   30,
			Description: `Deskripsi Samsung 15.6" Notebook 5 NT500R5W-XD7S i7-7500U 2.7GHz DDR4 8GB SSD 25
			Samsung 15.6" Notebook 5 NT500R5W-XD7S i7-7500U 2.7GHz DDR4 8GB SSD 256GB 940MX.
			`,
			Price: 23900000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, samsungNfNc)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-nb-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-nb-2.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-nb-3.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-nb-4.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-nb-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 4
	go func() {
		defer seedProductInMerchantMike.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		dellLatitude := adapter.ProductCreateInput{
			Name: "DELL Latitude 7490 Non Touch (Core i7-8650U)",
			Category: domain.Category{
				Top:       "komputer & laptop",
				SecondSub: "komputer & laptop",
				ThirdSub:  "laptop",
			},
			Etalase: "laptop",
			Tags: []string{
				"laptop bagus",
				"laptop dell",
				"laptop ringan",
				"dell notebook",
				"laptop terbaru",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  30,
			Height: 30,
			Long:   30,
			Description: `Processor: Intel Core i7-8650U
			RAM: 8GB DDR4
			SSD: 512GB
			VGA: Intel UHD Graphics
			Konektivitas: Wifi + Bluetooth
			Ukuran Layar: 14 Inch
			HD
			Sistem Operasi: Windows 10 Pro
			Aksesoris: Dell HDMI to VGA Adapter - Oxti + Kit-Dell Essential Backpack 15 - S&P.
			`,
			Price: 26636000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, dellLatitude)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "dell-79-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "dell-79-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "dell-79-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "dell-79-4.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "dell-79-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 5
	go func() {
		defer seedProductInMerchantMike.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		acerLaptop := adapter.ProductCreateInput{
			Name: "ASUS Business Notebook ExpertBook P2451FB-EK7850R",
			Category: domain.Category{
				Top:       "komputer & laptop",
				SecondSub: "komputer & laptop",
				ThirdSub:  "laptop",
			},
			Etalase: "laptop",
			Tags: []string{
				"laptop bagus",
				"laptop acer",
				"laptop ringan",
				"acer notebook",
				"laptop terbaru",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  30,
			Height: 30,
			Long:   30,
			Description: `Processor: Intel Core i7-10510
			RAM: 8GB DDR4
			SSD: 512GB
			VGA: NVIDIA GeForce MX110
			Ukuran Layar: 14 Inch FHD
			Security: Fingerprint
			Konektivitas: Bluetooth + Wifi + LAN
			Sistem Operasi: Windows 10 Pro
			`,
			Price: 21044000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, acerLaptop)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "acer-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "acer-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "acer-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "acer-4.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "acer-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 6
	go func() {
		defer seedProductInMerchantMike.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		asus := adapter.ProductCreateInput{
			Name: "MSI GE63 Raider 8RE RGB Edition (GeForce GTX 1060 6GB) 9S7-16P512-239",
			Category: domain.Category{
				Top:       "komputer & laptop",
				SecondSub: "komputer & laptop",
				ThirdSub:  "laptop",
			},
			Etalase: "laptop",
			Tags: []string{
				"laptop bagus",
				"notebook bagus",
				"laptop asus",
				"laptop ringan",
				"asus notebook",
				"laptop terbaru",
				"asus laptop",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  30,
			Height: 30,
			Long:   30,
			Description: `Processor: Intel Core i7-8750H
			RAM: 8GB x 2 DDR4
			HDD: 1TB
			SSD: 256GB
			VGA: GeForce GTX 1060 6GB
			Ukuran layar: 15.6 Inch FHD
			Sistem Operasi: Windows 10 Home
			`,
			Price: 24999000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, asus)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "msi-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "msi-2.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "msi-3.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "msi-4.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "msi-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 7
	go func() {
		defer seedProductInMerchantMike.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		asus := adapter.ProductCreateInput{
			Name: "MSI GE63 Raider 8RE RGB Edition (GeForce GTX 1060 6GB) 9S7-16P512-239",
			Category: domain.Category{
				Top:       "komputer & laptop",
				SecondSub: "komputer & laptop",
				ThirdSub:  "laptop",
			},
			Etalase: "laptop",
			Tags: []string{
				"laptop bagus",
				"notebook bagus",
				"laptop asus",
				"laptop ringan",
				"asus notebook",
				"laptop terbaru",
				"asus laptop",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  30,
			Height: 30,
			Long:   30,
			Description: `Processor: Intel Core i7-8750H
			RAM: 8GB x 2 DDR4
			HDD: 1TB
			SSD: 256GB
			VGA: GeForce GTX 1060 6GB
			Ukuran layar: 15.6 Inch FHD
			Sistem Operasi: Windows 10 Home
			`,
			Price: 24999000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, asus)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "msi-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "msi-2.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "msi-3.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "msi-4.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "msi-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 8
	go func() {
		defer seedProductInMerchantMike.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		asus := adapter.ProductCreateInput{
			Name: "ACER Swift 3 AMD Ryzen Non Windows (SF315-41) [NX.GV7SN.002/UN.GV7SN.001] - Silver",
			Category: domain.Category{
				Top:       "komputer & laptop",
				SecondSub: "komputer & laptop",
				ThirdSub:  "laptop",
			},
			Etalase: "laptop",
			Tags: []string{
				"laptop bagus",
				"notebook bagus",
				"laptop acer",
				"laptop ringan",
				"acer notebook",
				"laptop terbaru",
				"acer laptop",
				"acer swift",
				"acer swift 3",
				"laptop pelajar",
				"laptop ringan",
				"laptop sekolah",
				"laptop bisnis",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  30,
			Height: 30,
			Long:   30,
			Description: `Processor: AMD Ryzen 7 2700U
			RAM: 8GB DDR4
			HDD: 1TB
			Ukuran Layar: 15.6 Inch FHD
			VGA: Radeon Vega 10 Graphics
			Sistem Operasi: Non OS
			`,
			Price: 13005000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, asus)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "swift-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "swift-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "swift-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "swift-4.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "swift-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 9
	go func() {
		defer seedProductInMerchantMike.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		asus := adapter.ProductCreateInput{
			Name: "APPLE MacBook Pro [MV932ID/A] - Silver",
			Category: domain.Category{
				Top:       "komputer & laptop",
				SecondSub: "komputer & laptop",
				ThirdSub:  "laptop",
			},
			Etalase: "laptop",
			Tags: []string{
				"laptop bagus",
				"notebook bagus",
				"laptop apple",
				"laptop ringan",
				"apple notebook",
				"laptop terbaru",
				"apple laptop",
				"apple pro",
				"apple macbook",
				"laptop pelajar",
				"laptop ringan",
				"laptop sekolah",
				"laptop bisnis",
				"macbook pro",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  30,
			Height: 30,
			Long:   30,
			Description: `Processor: Intel Core i9 2.3 GHz
			RAM: 16GB DDR4
			SSD: 512GB
			VGA: Radeon Pro 560X with 4GB
			Konektivitas: Wifi + Bluetooth
			Ukuran Layar: 15.4 Inch
			Sistem Operasi: Mac OS
			`,
			Price: 39831000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, asus)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "apple-pro-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "apple-pro-2.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "apple-pro-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "apple-pro-4.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "apple-pro-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 10
	go func() {
		defer seedProductInMerchantMike.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		asus := adapter.ProductCreateInput{
			Name: "APPLE MacBook Air [MVFN2ID/A] - 256GB - Gold",
			Category: domain.Category{
				Top:       "komputer & laptop",
				SecondSub: "komputer & laptop",
				ThirdSub:  "laptop",
			},
			Etalase: "laptop",
			Tags: []string{
				"laptop bagus",
				"notebook bagus",
				"laptop apple",
				"laptop ringan",
				"apple notebook",
				"laptop terbaru",
				"apple laptop",
				"apple pro",
				"apple macbook",
				"laptop pelajar",
				"laptop ringan",
				"laptop sekolah",
				"laptop bisnis",
				"macbook air",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  30,
			Height: 30,
			Long:   30,
			Description: `Processor: Intel Core i5 8th 1.6GHz
			RAM: 8GB LPDDR3
			Grafik: Intel UHD Graphics
			Konektivitas: Wifi + Bluetooth
			Camera
			Ukuran Layar: 13.3 Inch
			Sistem Operasi: MacOS
			`,
			Price: 16490000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, asus)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "apple-air-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "apple-air-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "apple-air-3.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "apple-air-4.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "apple-air-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Mike: failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	seedProductInMerchantMike.Wait()
	if len(errs) != 0 {
		return products, errs[0]
	}

	log.Println("Seed Product Merchant Mike: finish!")
	return products, nil
}

func seedProductInJainalMerchant(productUsecase logic.ProductUsecase, credential domain.Credential) ([]domain.Product, error) {
	log.SetOutput(os.Stdout)
	log.Println("Seed Product Merchant Jainal starting!")

	var products []domain.Product
	var errs []error

	var seedProductInMerchantJainal sync.WaitGroup
	seedProductInMerchantJainal.Add(10)

	//product 1
	go func() {
		defer seedProductInMerchantJainal.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		samsungGalaxyFlex := adapter.ProductCreateInput{
			Name: "SAMSUNG Galaxy Note20 8GB/256GB - Mystic Green",
			Category: domain.Category{
				Top:       "handphone & tablet",
				SecondSub: "handphone & tablet",
				ThirdSub:  "handpone",
			},
			Etalase: "handphone",
			Tags: []string{
				"hp samsung",
				"samsung galaxy",
				"samsung 5G",
				"handphone 5G",
				"handphone 8GB",
				"samsung note",
			},
			Colors: []string{"hijau"},
			Weight: 1000,
			Width:  10,
			Height: 10,
			Long:   10,
			Description: `Chipset: Octa-Core 2.73GHz + 2.6GHz + 2GHz
			Kamera Belakang: 12MP & 64MP & 12MP
			Kamera Depan: 10MP
			Ukuran Layar: 6.7 Inch Super AMOLED
			Battery: 4300 mAh
			Sistem Operasi: Android 10`,
			Price: 13599000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, samsungGalaxyFlex)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-note-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-note-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-note-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-note-4.png"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-note-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 2
	go func() {
		defer seedProductInMerchantJainal.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		samsungNfNc := adapter.ProductCreateInput{
			Name: "SAMSUNG Galaxy S20 FE 8GB/256GB - Cloud Mint",
			Category: domain.Category{
				Top:       "handphone & tablet",
				SecondSub: "handphone & tablet",
				ThirdSub:  "handpone",
			},
			Etalase: "handphone",
			Tags: []string{
				"hp samsung",
				"samsung galaxy",
				"samsung 4G",
				"handphone 8GB",
				"samsung note",
				"handphone 4G",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  10,
			Height: 10,
			Long:   10,
			Description: `Chipset: Exynos 990
			Kamera Belakang: 12MP & 8MP & 12MP
			Kamera Depan: 32MP
			Ukuran Layar: 6.5 Inch 120Hz
			Battery: 4500 mAh
			IP68
			Sistem Operasi: Android 10.0; One UI 2
			`,
			Price: 10999000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, samsungNfNc)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-s20-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-s20-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-s20-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-s20-4.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "samsung-s20-5.jpeg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 3
	go func() {
		defer seedProductInMerchantJainal.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		samsungNfNc := adapter.ProductCreateInput{
			Name: "Realme C17 6GB/256GB - Lake Green",
			Category: domain.Category{
				Top:       "handphone & tablet",
				SecondSub: "handphone & tablet",
				ThirdSub:  "handpone",
			},
			Etalase: "handphone",
			Tags: []string{
				"hp realme",
				"handphone realme",
				"samsung 4G",
				"handphone 6GB",
				"realme green",
				"handphone 4G",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  10,
			Height: 10,
			Long:   10,
			Description: `Chipset: Snapdragon 460
			Kamera Belakang: 13MP main camera + 8MP ultra wide-angle + 2MP B&W camera + 2MP macro camera
			Kamera depan: 8MP In-display Selfie
			Ukuran Layar: 6.5 Inch 90Hz
			Security: Fingerprint
			Battery: 5000mAh
			Sistem Operasi: Android 10 realme UI
			`,
			Price: 2849000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, samsungNfNc)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "realme-c17-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "realme-c17-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "realme-c17-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "realme-c17-4.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "realme-c17-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 4
	go func() {
		defer seedProductInMerchantJainal.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		dellLatitude := adapter.ProductCreateInput{
			Name: "XIAOMI Poco X3 NFC 8GB/128GB - Grey",
			Category: domain.Category{
				Top:       "handphone & tablet",
				SecondSub: "handphone & tablet",
				ThirdSub:  "handpone",
			},
			Etalase: "handphone",
			Tags: []string{
				"hp xiaomi",
				"handphone xiaomi",
				"samsung 4G",
				"handphone 8GB",
				"handphone 4G",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  10,
			Height: 10,
			Long:   10,
			Description: `Chipset: Qualcomm SM7150-AC Snapdragon 732G (8 nm)
			Kamera Belakang: 64 MP + 13 MP + 2 MP + 2 MP
			Kamera Depan: 20 MP
			Ukuran Layar: 6.67 inch
			Baterai: 5160 mAh
			Sistem Operasi: Android 10
			`,
			Price: 3999000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, dellLatitude)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "xiaomi-poco-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "xiaomi-poco-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "xiaomi-poco-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "xiaomi-poco-4.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "xiaomi-poco-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 5
	go func() {
		defer seedProductInMerchantJainal.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		acerLaptop := adapter.ProductCreateInput{
			Name: "HUAWEI P40 8GB/128GB - Silver Frost",
			Category: domain.Category{
				Top:       "handphone & tablet",
				SecondSub: "handphone & tablet",
				ThirdSub:  "handpone",
			},
			Etalase: "handphone",
			Tags: []string{
				"hp huawei",
				"handphone huawei",
				"huawei 4G",
				"handphone 8GB",
				"handphone 4G",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  10,
			Height: 10,
			Long:   10,
			Description: `Processor: Kirin 990 5G
			Kamera belakang: 50 MP + 16 MP
			Kamera depan: 32 MP
			Ukuran layar: 6.1 inch
			Resolusi: FHD+ 2340 x 1080 Pixel
			Baterai: 3800 mAh
			Sistem operasi: EMUI 10.1 (Based on Android 10)
			`,
			Price: 9899000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, acerLaptop)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "huawei-p40-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "huawei-p40-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "huawei-p40-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "huawei-p40-4.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "huawei-p40-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 6
	go func() {
		defer seedProductInMerchantJainal.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		asus := adapter.ProductCreateInput{
			Name: "VIVO Y12i 3GB/32GB - Mineral Blue",
			Category: domain.Category{
				Top:       "handphone & tablet",
				SecondSub: "handphone & tablet",
				ThirdSub:  "handpone",
			},
			Etalase: "handphone",
			Tags: []string{
				"hp vivo",
				"handphone vivo",
				"vivo 4G",
				"handphone 8GB",
				"handphone 4G",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  10,
			Height: 10,
			Long:   10,
			Description: `Prosesor: Snapdragon 439
			Layar: 6.35 inci
			Kamera: Depan 8MP
			Kamera Belakang: 13MP & 2MP
			Baterai: 5000mAh (TYP)
			Sistem Operasi: Android 9.0 Pie
			`,
			Price: 1899000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, asus)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "vivo-y12-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "vivo-y12-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "vivo-y12-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "vivo-y12-4.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "vivo-y12-5.jpeg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 7
	go func() {
		defer seedProductInMerchantJainal.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		asus := adapter.ProductCreateInput{
			Name: "NOKIA 6.1 Plus 4GB/64GB - Black",
			Category: domain.Category{
				Top:       "handphone & tablet",
				SecondSub: "handphone & tablet",
				ThirdSub:  "handpone",
			},
			Etalase: "handphone",
			Tags: []string{
				"hp nokia",
				"handphone nokia",
				"nokia 4G",
				"handphone 8GB",
				"handphone 4G",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  10,
			Height: 10,
			Long:   10,
			Description: `Processor: Qualcomm SDM636 Snapdragon 636 Octa-core 1.8 GHz Kryo 260
			Ukuran Layar: 5.8 inch
			Kamera Belakang: 16 MP+ 5 MP
			Kamera Depan: 16 MP
			Baterai: 3060 mAh
			Dual SIM
			Fingerprint
			Android OS: Android 8.1 Oreo
			`,
			Price: 3399000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, asus)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "nokia-6-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "nokia-6-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "nokia-6-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "nokia-6-4.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "nokia-6-5.jpeg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 8
	go func() {
		defer seedProductInMerchantJainal.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		asus := adapter.ProductCreateInput{
			Name: "OPPO Reno3 Pro 8GB/256GB - Midnight Black",
			Category: domain.Category{
				Top:       "handphone & tablet",
				SecondSub: "handphone & tablet",
				ThirdSub:  "handpone",
			},
			Etalase: "handphone",
			Tags: []string{
				"hp oppo",
				"handphone oppo",
				"oppo 4G",
				"handphone 8GB",
				"handphone 4G",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  10,
			Height: 10,
			Long:   10,
			Description: `Chipset: Mediatek Helio P95
			Kamera Belakang: 64MP + 13MP + 8MP + 2MP
			Kamera Depan: 44MP + 2MP
			Ukuran Layar: 6.4 Inch
			Battery: 4025 mAh
			Operating System: ColorOS 7 Based on Android 10
			`,
			Price: 7299000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, asus)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "oppo-reno-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "oppo-reno-2.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "oppo-reno-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "oppo-reno-4.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "oppo-reno-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 9
	go func() {
		defer seedProductInMerchantJainal.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		asus := adapter.ProductCreateInput{
			Name: "SHARP Aquos Zero 2 8GB/256GB - Black",
			Category: domain.Category{
				Top:       "handphone & tablet",
				SecondSub: "handphone & tablet",
				ThirdSub:  "handpone",
			},
			Etalase: "handphone",
			Tags: []string{
				"hp sharp",
				"handphone sharp",
				"sharp 4G",
				"handphone 8GB",
				"handphone 4G",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  10,
			Height: 10,
			Long:   10,
			Description: `Processor: Intel Core i9 2.3 GHz
			RAM: 16GB DDR4
			SSD: 512GB
			VGA: Radeon Pro 560X with 4GB
			Konektivitas: Wifi + Bluetooth
			Ukuran Layar: 15.4 Inch
			Sistem Operasi: Mac OS
			`,
			Price: 12999000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, asus)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "sharp-aquos-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "sharp-aquos-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "sharp-aquos-3.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "sharp-aquos-4.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "sharp-aquos-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	//product 10
	go func() {
		defer seedProductInMerchantJainal.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ctx = context.WithValue(ctx, "credential", credential)
		defer cancel()

		asus := adapter.ProductCreateInput{
			Name: "ASUS ROG Phone 3 12GB/256GB - Black Glare",
			Category: domain.Category{
				Top:       "handphone & tablet",
				SecondSub: "handphone & tablet",
				ThirdSub:  "handpone",
			},
			Etalase: "handphone",
			Tags: []string{
				"hp asus",
				"asus rog",
				"handphone asus",
				"asus 4G",
				"handphone 8GB",
				"handphone 4G",
				"handphone asus rog",
			},
			Colors: []string{"biru"},
			Weight: 1000,
			Width:  10,
			Height: 10,
			Long:   10,
			Description: `Chipset: Qualcomm SM8250 Snapdragon 865 Plus
			Kemera Belakang: 64MP + 13MP + 5MP
			Kamera Depan: 24MP
			Ukuran Layar: 6.59 Inch
			Baterai: 6000mAh
			Sistem Operasi: Android 10
			`,
			Price: 14999000,
			Stock: 10,
		}

		product, err := productUsecase.Create(ctx, asus)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}

		//give distance saving product in elastic
		time.Sleep(30 * time.Millisecond)
		ctxUploadPhotos, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		photos := []string{
			fmt.Sprintf("%s%s", googleStorageAPI, "asus-rog-1.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "asus-rog-2.jpg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "asus-rog-3.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "asus-rog-4.jpeg"),
			fmt.Sprintf("%s%s", googleStorageAPI, "asus-rog-5.jpg"),
		}
		product, err = productUsecase.UploadPhotos(ctxUploadPhotos, photos, product.ID)
		if err != nil {
			log.Printf("Seed Product Merchant Jainal failed, cause %s \n", err)
			errs = append(errs, err)
			return
		}
		products = append(products, product)
	}()

	seedProductInMerchantJainal.Wait()
	if len(errs) != 0 {
		return products, errs[0]
	}

	log.Println("Seed Product Merchant Jainal finish!")
	return products, nil
}
