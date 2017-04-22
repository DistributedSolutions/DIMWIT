package torrent

import (
	"fmt"
	"github.com/anacrolix/torrent"
)

//The point of this is to separate the Torrent Client from the body of access info from it
type ClientInterface struct {
	torrentClient *TorrentClient
}

func (c *ClientInterface) SetClient(client *TorrentClient) {
	c.torrentClient = client
}

type JSONFiles struct {
	Files []*JSONFile `json:"files"`
}

type JSONFile struct {
	Name       string                `json:"fileName"`
	Size       int64                 `json:"byteSize"`
	FilePieces []*JSONFilePieceState `json:"filePieceState"`
}

//COPY OF PIECE STATE FROM anacrolix
type JSONFilePieceState struct {
	Percentage float64         `json:"percentage"` // Bytes within the piece that are part of this File.
	PieceState *JSONPieceState `json:"pieceState"`
}

func ToJSONFilePieceState(percentage float64, filePieceState torrent.PieceState) *JSONFilePieceState {
	jfp := JSONFilePieceState{
		percentage,
		ToJSONPieceState(&filePieceState),
	}
	return &jfp
}

//COPY OF PIECE STATE FROM anacrolix
type JSONPieceState struct {
	// The piece is available in its entirety.
	Complete bool `json:"completed"`
	// The piece is being hashed, or is queued for hash.
	Checking bool `json:"checking"`
	// Some of the piece has been obtained.
	Partial bool `json:"partial"`
}

func ToJSONPieceState(filePieceState *torrent.PieceState) *JSONPieceState {
	jfps := JSONPieceState{
		filePieceState.Complete,
		filePieceState.Checking,
		filePieceState.Partial,
	}
	return &jfps
}

// func (c *ClientInterface) GetTorrentFileMetaData(torrentHash string) (*JSONFiles, error) {
// 	torrentFiles, err := c.torrentClient.GetTorrentFiles(torrentHash)
// 	if err != nil {
// 		return nil, err
// 	}
// 	jsonFiles := new(JSONFiles)
// 	jsonFiles.Files = make([]*JSONFile, len(torrentFiles), len(torrentFiles))
// 	// each torrent file
// 	for tfi, torrentFile := range torrentFiles {

// 		torrentFilePiecesStates := torrentFile.State()
// 		jsonFile := new(JSONFile)
// 		jsonFile.Size = torrentFile.Length()
// 		jsonFile.Name = torrentFile.Path()
// 		jsonFile.FilePieces = make([]*JSONFilePieceState, len(torrentFilePiecesStates), len(torrentFilePiecesStates))

// 		// each torrent file pieces
// 		for tfpsi, tfps := range torrentFilePiecesStates {
// 			jsonFile.FilePieces[tfpsi] = ToJSONFilePieceState(tfps)
// 		}
// 		jsonFiles.Files[tfi] = jsonFile
// 	}
// 	return jsonFiles, nil
// }

//chuncks the sizes up into meaning pieces
func (c *ClientInterface) GetTorrentFileMetaDataChunked(torrentHash string) (*JSONFiles, error) {
	torrentFiles, err := c.torrentClient.GetTorrentFiles(torrentHash)
	if err != nil {
		return nil, err
	}
	jsonFiles := new(JSONFiles)
	jsonFiles.Files = make([]*JSONFile, len(torrentFiles), len(torrentFiles))
	// each torrent file
	for tfi, torrentFile := range torrentFiles {

		torrentFilePiecesStates := torrentFile.State()
		jsonFile := new(JSONFile)
		jsonFile.Size = torrentFile.Length()
		jsonFile.Name = torrentFile.Path()
		// jsonFile.FilePieces = make([]*JSONFilePieceState, len(torrentFilePiecesStates), len(torrentFilePiecesStates))
		var jsonChunk []*JSONFilePieceState

		var totalByteCount int64
		totalByteCount = 0
		var byteSize int64
		byteSize = 0
		var tempPerc float64
		var totalTempPerc float64
		totalTempPerc = 0
		// each torrent file pieces
		for tfpsi, _ := range torrentFilePiecesStates {
			if tfpsi == len(torrentFilePiecesStates)-1 {
				//its the last index
				byteSize += torrentFilePiecesStates[tfpsi].Bytes
				tempPerc = 100 * (float64(byteSize) / float64(jsonFile.Size))
				jsonChunk = append(jsonChunk, ToJSONFilePieceState(tempPerc, torrentFilePiecesStates[tfpsi].PieceState))
				totalTempPerc += tempPerc
			} else if tfpsi > 0 && (torrentFilePiecesStates[tfpsi].Complete != torrentFilePiecesStates[tfpsi-1].Complete ||
				torrentFilePiecesStates[tfpsi].Checking != torrentFilePiecesStates[tfpsi-1].Checking ||
				torrentFilePiecesStates[tfpsi].Partial != torrentFilePiecesStates[tfpsi-1].Partial) {
				//not all equal therefore add and then start over
				tempPerc = 100 * (float64(byteSize) / float64(jsonFile.Size))
				jsonChunk = append(jsonChunk, ToJSONFilePieceState(tempPerc, torrentFilePiecesStates[tfpsi].PieceState))
				byteSize = torrentFilePiecesStates[tfpsi].Bytes
				totalTempPerc += tempPerc
			} else {
				byteSize += torrentFilePiecesStates[tfpsi].Bytes
			}
			totalByteCount += torrentFilePiecesStates[tfpsi].Bytes
		}
		fmt.Printf("Total Count: %d and this count: %d percentage: %.6f\n", totalByteCount, torrentFile.Length(), totalTempPerc)
		jsonFile.FilePieces = jsonChunk
		jsonFiles.Files[tfi] = jsonFile
	}
	return jsonFiles, nil
}
