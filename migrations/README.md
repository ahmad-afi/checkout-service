# hera - database migration

# Cara Menggunkan golang migration

1. Install dahulu [disini](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md)
2. Buat file migrate dengan command `make migrate name=create-table-user`
3. Ubah isi file up dan down, cek di folder `migrations`
4. Untuk migration semua file gunakan command `make migrate-up` jika ingin bebrapa step `make migrate-up step=10`. Nilai default adalah 1.
5. Untuk rollback semua file gunakan command `make migrate-down` jika ingin bebrapa step `make migrate-down step=10`. Nilai default adalah 1.

# PENTING !!!!

1. Jangan mengedit file migrations yang lama atau sudah di up. karena perubahan query tidak akan dijalankan.
2. Jika ingin running seluruh migration harap gunakan  
   `make migrate-up steps=0 forceMigration=true`

# HARAP BIJAK DALAM MENJALANKAN MIGRATION

# KESALAHAN DATA ADALAH TANGGUNG JAWAB ANDA

# Jika terdapat dirty state

1. Misal nama file terakhir Anda adalah `20231218110736_create-table-users.down.sql`, maka lihat lah nama file sebelumnya `20231218101715_create-table-partner.down.sql` lalu ambil timestamp nya
2. Gunakan perintah `make migrate-force version={version}` contoh `make migrate-force version=1`.
   Sesuaikan dengan timestamp file ke dua terakhir Anda.
3. Lalu lakukan perubahan query di file migration terbaru.
4. Jalankan lagi `make migrate-up`
