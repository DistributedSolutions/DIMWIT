# Intro
The use of Factom is purely a storage solution. The users gather all the data and arrange it locally for quicker searches. Because of this, the actual arrangement is Factom is not critical to performance, but purely for organizational purposes when downloading everything needed.

Prefix is 0xdcff

# Dictionary

Short dictionary on the vocab being used throughout

|Term|Definition|
|---|---|
|Client|The TorrentTube application. Downloads and caches data from Factom, capable of streaming content, and (not required) seeding content. Users may wish to use alternative torrent clients for seeding.|
|User|A person interacting with a client. Users can only stream and seed.|
|Channel|A user or collection of users that upload content. Channels can only upload and seed content|
|Content|Any data able to be torrented, usually in the form of a video.|
|Torrent|A torrent is the .torrent file metadata needed to successfully obtain content|
|Peers|Specifically clients that are seeding the requested torrent. These clients could also be using other applications to seed, but when reffered to as 'peers', we do not care what they are using as their application.|
|Tracker|Where to direct the client to find peers with the requested torrent|
|Factom Data|This data is downloaded from the Factom Blockchain. It is used to organize all the data needed to find content|
|Factom Entries|1Kb to 10Kb chunks of data entered into factom|
|Factom Chain|A chain is a linked list of entries|

# Info needed to be stored on Factom
There isn't many data types to be kept track of. It comes down to Channels and their content. Everything else is handled by the clients off the Factom Blockchain
- Torrents
  - Torrent metadata
  - Content metadata (we want additional metadata about the content alongside torrent metadata)
- Channels - We need to be able to identify where content comes from.
  - Channel Identity (Signing keys)
  - Channel Content


# Factom Chains
## Each Chain has:

- Header Entry
  - The first entry in a chain. This usually sets up public keys that will be used to sign subsequent entries. It is impossible to prevent spam, but using signatures allows spam to be identified and ignored.
- Entries
  - Entries following the first entry
  - Every entry (including the header entry) have a header section themselves in the form of "External IDs". External IDs are a list of values, the list is of any size.
    - Used to identify spam entries
  - Every entry has a content section, which is optional and can contain anything.

## Factom Data Section
All entries entered into factom will fall in one of these sections
- Master Section
 - Only 1 master chain exists
- Channel Sections
- Content Sections

When reffering to a "ChainID" that means there is a 32 byte hash that can be used to find the chain within factom. Entries also have a hash that can be used to find an individual entry, but most entries are found by traversing chains.

All entries will have a version number. If we change formatting later, we will need to add legacy support. Knowing the version number will allow that

## Master Section
There is 1 chain, in which all clients have a hardcoded link to. This is the chain that will have links/directions capable of finding all content uploaded.

|Header Entry||
|---|---|
|ExtID (0)|Version|
|ExtID (1)|"Master Chain"|
|ExtID (2)|nonce|
|ExtID (3)|Unsure|
|Content | Unsure|

Entries will be made by channels to register themselves in the master list.

|Entries||
|---|---|
|ExtID (0)|Version|
|ExtID (1)|"Channel Chain"|
|ExtID (2)|Channel Root ChainID|
|ExtID (3)|Public Key (3)|
|ExtID (4)|Signature of ExtID(0-3)|
|Content|Unsure|

## Channel Sections
Consists of 3 Chains
- Channel Root Chain
  - The main chain that has links to the channel management & content chains
  - Used for changing keys being used to sign
- Channel Management Chain
  - Used to change channel metadata. Changes here affect what users see
    - E.g: Channel descriptions and such
- Content Chain
  - A chain for linking to specific content chains

### Channel Root Chain
The reason for 3 keys, is if the level 3 key is compromised, the level 2 key can change the level 3. If level 2 is compromised, level 1 can change level 2. If level 1 is compromised, you just lost your chain bro. The nonce is so we can make the chains start with some specific hex value.

|Header Entry||
|---|---|
|ExtID (0)|Version {1 byte}|
|ExtID (1)|"Channel Root Chain" {18 bytes}|
|ExtID (2)|PublicKey 1 {32 bytes}|
|ExtID (3)|PublicKey 2 {32 bytes}|
|ExtID (4)|PublicKey 3 {32 bytes}|
|ExtID (5)|nonce {8 bytes}|
|Content | Unsure (possible metadata?)|

Entry to designate a Content signing key. This can be changed by the (3) public key

|Entry||
|---|---|
|ExtID (0)|Version {1 byte}|
|ExtID (1)|"Content Signing Key" {19 byte}|
|ExtID (2)|Content Signing Key {32 bytes}|
|ExtID (3)|Timestamp {8 bytes}|
|ExtID (4)|Signature of ExtID(0-3) {64 bytes}|
|Content|Unsure|

- Timestamp must be withing +/- 6 hrs of the block it is in to prevent replays


Entries will also be made by the channel to register there 2 other chains

|Entry||
|---|---|
|ExtID (0)|Version {1 byte}|
|ExtID (1)|"Register Management Chain" {25 bytes}|
|ExtID (2)|Channel Management ChainID {32 bytes}|
|ExtID (3)|Public Key (3) {32 bytes}|
|ExtID (4)|Signature of ExtID(0-2) {64 bytes}|
|Content|Unsure|

|Entry||
|---|---|
|ExtID (0)|Version {1 byte}|
|ExtID (1)|"Register Content Chain" {22 byte}|
|ExtID (2)|Channel Content ChainID {32 bytes}|
|ExtID (3)|Public Key (3) {32 bytes}|
|ExtID (4)|Signature of ExtID(0-2) {64 bytes}|
|Content|Unsure|

### Channel Manage Chain

|Header Entry||
|---|---|
|ExtID (0)|Version {1 byte}|
|ExtID (1)|"Channel Management Chain" {24 bytes}|
|ExtID (2)|Channel Root ChainID {32 bytes}|
|ExtID (3)|Public Key (3) {32 bytes}|
|ExtID (4)|Signature of ExtID(0-2) {64 bytes}|
|ExtID (5)|nonce {8 bytes}|
|Content | Channel Title|

Entries unknown atm. All metadata changes.

### Channel Content Chain
|Header Entry||
|---|---|
|ExtID (0)|Version {1 byte}|
|ExtID (1)|"Channel Content Chain" {21 bytes}|
|ExtID (2)|Channel Root ChainID {32 bytes}|
|ExtID (3)|Public Key (3) {32 bytes}|
|ExtID (4)|Signature of ExtID(0-2) {64 bytes}|
|ExtId (5)| nonce {8 bytes}|
|Content | Unsure|

Entries will point to individual content chains

|Entry||
|---|---|
|ExtID (0)|Version {1 byte}|
|ExtID (1)|Content Type {1 byte}|
|ExtID (2)|"Content Link" {12 bytes}|
|ExtID (3)|Channel Root ChainID {32 bytes}|
|ExtID (4)|Timestamp {8 bytes}|
|ExtID (5)|Content Signing Key {32 bytes}|
|ExtID (6)|Signature of ExtID(0-4) {64 bytes}|
|Content|Unsure|

The content type is a byte indicator of the type of content. If we support more types of content in the future, this will be important, but for now, its mostly all of type "video", or 0x00.

We don't necessarily need to timestamp, but if we do not, then someone could replay the same content over and over again. Clients should be smart enough to recognize duplicates, but the timestamp makes it easier to ignore. +/- 6 hrs from block (may change).

### Content Chain
A new chain will be made per piece of content. This allows us to add more functionality in the future, and possibly even commenting.

|Header Entry||
|---|---|
|ExtID (0)|Version {1 byte}|
|ExtID (1)|Content Type {1 byte}|
|ExtID (2)|Total Entries {1 byte}|
|ExtID (3)|"Content Chain" {13 bytes}|
|ExtID (4)|Channel Root ChainID {32 bytes}|
|ExtID (5)|InfoHash {20 bytes}|
|ExtID (6)|Timestamp {8 bytes}|
|ExtID (7)|Shift Cipher Key {1 byte}|
|ExtID (8)|Content Signing Key {32 bytes}|
|ExtID (9)|Signature of ExtID(0-7) {64 bytes}|
|ExtID (10)|nonce {8 bytes}|
|Content|Shift Ciphered Content Metadata|

This has quite a bit of data. The Channel Root ChainID is included to allow backwards traversing. So if someone links you to a video, the client can backwards traverse for the channel keys.

The infohash is the Torrent infohash. It will be redeclared in the body, but this is the unique identifier for torrent and should be included in the signature.

The timestamp is not to prevent replays, we don't care about chain replays. The reason why it is there, is because if someone links you to a video, they are linking straight to a chain, and instead of having to do a lookup, we can just take the timestamp from the entry. !NOTE!: If possible, do not trust the timestamp for anything regarding which content signing key is the current valid. Use the entryblock timestamp if we have it. TBH, might remove the timestamp, still thinking about how well it can be trusted...

Encrypted Binary Block: **Might remove, depending on encoding** We don't want to upload plaintext data, so we encrypt the metadata with a quick encryption (have not decided which, honestly anything will work) to prevent anyone not using the system to immediately recognize the data and it's purpose. All data entered into Factom is public, and it can't hurt to obfuscate the torrent data.

All torrent metadata is in the content. Have not totally decided what we need here

|Encrypted Binary Block|
|---|---|
|Content Metadata|
|Title|
|Channel Root ChainID (we always want to grab this data from the source to prevent impersonation)|
|Description|
|Tags|
|---|
|Torrent Metadata|
|Everything in a .torrent file|