﻿SELECT RESERVED
FROM SEATS
INNER JOIN WAGON_SEAT_CONNECTION ON ID = SEATS_ID
INNER JOIN WAGONS ON WAGONS.ID = WAGON_SEAT_CONNECTION.WAGONS_ID
WHERE WAGONS.ID = "Vagon1" AND SEATS.ID = 1