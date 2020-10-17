# JK Hazard API

## Headers

Every request must contain the header key, Token and a valid token string provided at NewRoom or JoinRoom.

> Token : <valid_token_id>

> Content-Type: application/json

## Endpoints

>GET /api/rooms

>POST /api/rooms

>PUT /api/rooms/join

>GET /api/room/{roomTAG}

>GET /api/room/{roomTAG}/start

>PUT /api/room/{roomTAG}/send-card

>GET /api/room/{roomTAG}/judge-cards

>PUT /api/room/{roomTAG}/round-winner

>GET /api/room/{roomTAG}/heart-beat

>GET /api/room/{roomTAG}/host

>GET /api/room/{roomTAG}/judge

>GET /api/room/{roomTAG}/player

## Details

> POST /api/rooms : New Room

Request Body: 

    {
        "username": "yourname"
    }

Response Body:

    {
        "room": {
        },
        "token": "uuid-string"
    }

>PUT /api/rooms/join : Join Room

Request Body: 

    {
        "roomTag": "string",
        "username": "chava"
    }

Response Body:

    {
        "token": "uuid-string"
    }

>PUT /api/room/{roomTAG}/send-card : Send card to be judged

Request Body: 

    {
        "cardID": "227"
    }

Response Body:

    {
        "newCard": {
            "ID": 265,
            "Src": "out-265",
            "IsBlack": false
        }
    }

>PUT /api/room/{roomTAG}/round-winner : Judge sends winner

Request Body: 

    {
        "username": "anyname"
    }

Response Body:

    {
        "message": "Round winner received and updated"
    }
