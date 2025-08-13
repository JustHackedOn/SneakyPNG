SneakyPNG 🥷
This Go tool hides files inside PNG images using RC4 encryption and extracts them back. 🖼️🔒 Use embed.go to hide a file in a PNG and fetch.go to get it out. For learning only! 😈
What It Does

Embed: Hides any file in a PNG’s IDAT chunks with RC4 encryption. Output PNG looks normal.
Fetch: Finds hidden file using a CRC hash, decrypts, and saves it. Can run .exe on Windows.

Setup

Install Go (1.16+): golang.org
Clone repo: git clone https://github.com/yourusername/SneakyPNG.git
Build: 
go build embed.go → embed (or embed.exe)
go build fetch.go → fetch (or fetch.exe)



Usage: Embed (Hide File) 🕵️‍♂️
./embed -i <secret-file> -png <input-pic.png> -o <output-pic.png>


-i: File to hide (e.g., .txt, .exe).
-png: Input PNG.
-o: Output PNG name.

Example:
./embed -i secret.txt -png cat.png -o sneaky-cat.png

Save the MARKED_IDAT_HASH (e.g., 0xABC123) for fetching.
Usage: Fetch (Extract File) 🕵️‍♀️
./fetch -png <sneaky-cat.png> -o <secret.txt> -crc <0xYourHash> [-exec false]


-png: PNG with hidden file.
-o: Output file name.
-crc: Hash from embed (e.g., 0xABC123).
-exec: Run .exe after? (Default: true, Windows only).

Example:
./fetch -png sneaky-cat.png -o secret.txt -crc 0xABC123

Made by [Your Name]. ⭐
