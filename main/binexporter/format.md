# Format
The first 4 bytes `00 00 00 00` are used to specify how far to jump in to get the the data of the world  
The next 4 bytes `00 00 00 00` are used to specify how far to jump in to start reading the pallet  
The next 2 bytes `00 00` are used to specify how wide each chunk is  
The next 2 bytes `00 00` are used to specify how deep each chunk is  
The next 2 bytes `00 00` are used to specify how tall each chunk is
The next 2 bytes `00 00` are used to specify the length of the pallet
So a full header should could look something like  `00 00 00 FF 00 00 00 0F 00 10 00 10 01 5E 00 0A`
| Hex | Meaning |
|--|--|
| 00 00 00 FF | This means that the data for the world starts 255 bytes in 
| 00 00 00 0F | This means that the data for the pallet starts 15 bytes in
| 00 10 | This means that each chunk is 16 wide 
| 00 10 | This means that each chunk is 16 deep
| 01 5E | This means that each chunk is 350 tall
| 00 0A | This means that there are 10 types of blocks in the pallet data

The pallet data is stored in a string in the format of `namespace:identifier` eg. `minecraft:stone`  
The world data should have a single byte `00` at the start that specifies if it uses RLE (run length encoding)
|Hex| Meaning |
|--|--|
| 00 | Does not use RLE |
| 01 | Does use RLE

after that the world  data will get written in the format of XYZ it should flatten each slice down into a 1D stream 
then it should move onto the next slice. Each block has a 2 byte identifier which should act as a direct index into the pallet eg. `0A` signifies that the block at that position has the value of the 0A index of the pallet (the 11nth item)  
Each chunk has to have a chunk start flag `FF AA FF AA` and a chunk end flag `AA FF AA FF`.  
All data between these two flags should get treated as compressed chunk data  
It is important to remember each chunk gets compressed with zip before getting added to the file.  
Before you attempt to read any data from a chunk make sure you have decompressed it

# Writing to disk
When the world data gets wrote to disk it will get compressed into a zip to reduce the size on disk so make sure to decompress it first before trying to parse the data

# None RLE chunk data
If the chunk data you are reading doesn't have the RLE flag set that means it stores the type of every block in the chunk 1 at a time
# RLE chunk data
if the chunk data you are reading has the RLE flag set that means that for each vertical strip of blocks on the Y axis it try to encode the block type then the length it runs for like this `00 00 00 FF`  
|Hex| Meaning |
|--|--|
| 00 00 | This means that we are using the block at index 0 in the pallet |
| 00 FF | This means that we use this same block for 255 blocks vertically |
 

