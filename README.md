PNG Payload Hider 🎨🔒
Yo, what's good? 😎 This Go project is your sneaky sidekick to hide files inside PNG images like a pro hacker! 🥷 Use embed.go to stash any file (payload) into a PNG with RC4 encryption, and fetch.go to pull it out and decrypt it. It's like hiding your snacks in a secret cookie jar! 🍪 Abdul Ahad says: "Security? Just an illusion!" 😂
Super simple, colorful logs, and works offline. Perfect for fun experiments! 🚀
What’s It Do? 🤔

Embed: Hides your file in a PNG using IDAT chunks and RC4 crypto. New PNG looks normal but has your secret! 📦➡️🖼️
Fetch: Finds the hidden file with a special CRC hash, decrypts, and saves it. Can even run .exe files on Windows! 🖼️➡️📦💥
Bonus: Cool ASCII art, progress bars, and rainbow logs! 🌈

Warning: For learning only! Don’t be naughty, or you’ll get "just hacked on"! 😜
Setup in a Snap 🛠️

Get Go: Install from golang.org (1.16+ works great).
Clone Repo: git clone https://github.com/yourusername/png-payload-hider.git
Build Tools:
go build embed.go → Makes embed (or embed.exe).
go build fetch.go → Makes fetch (or fetch.exe).


Dependencies: Uses color, progressbar, etc. Run go mod tidy if needed.

How to Hide (Embed) 🕵️‍♂️
Run ./embed to hide your file!
Command:
./embed -i <secret-file> -png <input-pic.png> -o <output-pic.png>


-i: Your file to hide (like .exe, .txt, anything!).
-png: A legit PNG image.
-o: Name for the sneaky output PNG.

Example:
./embed -i secret.txt -png meme.png -o sneaky-meme.png

Copy the MARKED_IDAT_HASH (like 0xABC123) it gives you for fetching! Logs saved too. 📝
How to Get It Back (Fetch) 🕵️‍♀️
Run ./fetch to extract!
Command:
./fetch -png <sneaky-meme.png> -o <secret.txt> -crc <0xYourHash> [-exec false] [-log mylog.txt]


-png: PNG with hidden stuff.
-o: Where to save the file.
-crc: That hash from embed (e.g., 0xABC123).
-exec: True/false—run .exe after? (Default: true, Windows only).
-log: Custom log file (default: timestamped).

Example:
./fetch -png sneaky-meme.png -o secret.txt -crc 0xABC123

Decrypted file pops out, and it runs if it’s an .exe and you said so! 📊
Pro Tips 💡

File Size: Big files make PNGs chonky—test with small ones first! 🐘
PNG Check: Auto-verifies if your PNG is real. No fakes! ❌
Fancy Logs: Colorful terminal + log file (red for errors, yellow for wins). 🌈
Play Safe: Try hiding a tiny .txt in a meme PNG for laughs! 🤣

Why It’s Awesome? 🌟

Easy-Peasy: Just run and hide. No PhD needed! 🏃‍♂️
Learn Cool Stuff: Steganography + RC4 = hacker vibes. 📈
Fun Vibes: Emojis and colors make coding a party! 🥳

Bugs or ideas? Open an issue or PR. Let’s make it sneakier! 😘
Made with ❤️ by [Your Name]. Star it if you love it! ⭐
