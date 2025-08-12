# escposimg

Go Library und cli-Tool, um Bilder mit ESC/POS-fähigen Belegdruckern auszudrucken.


## Hintergrund / Motivation

Viele Belegdrucker können/müssen mit dem Steuerprotokoll ESC/POS von Epson angesteuert werden.  Im Gegensatz zu den meisten modernen Bürodrucker, wird hierbei kein komplettes Dokumentlayout (etwa als PDF oder Post Script) geschickt, sondern einzelne Steuerbefehle wie »Text ausgeben« oder »Papier abschneiden«. Das eigentliche Rendering geschieht also auf dem Device. Das führt dazu, dass für das Erstellen von Layouts nicht die üblichen Tools (Typst, Word, LaTeX, Affinity Publisher etc) eingesetzt werden können.

Ein möglicher Workaround ist es, das gewünschte Dokument in ein monochromes (schwarz-weiß) Bild zu rasterizen und so an den Drucker zu schicken. Dies tut diese Library. Da Thermodrucker nur Schwarz drucken können, müssen Bilddateien, die Farben oder Graustufen enthalten, zunächst in ein Monochromes umgerechnet werden. Hierfür stellt das Tool eine Reihe unterschiedlicher Dithering-Algorithmen bereit. Die Druckdaten können über das Netzwerk an den Drucker geschickt werden.

Für alle gängigen Plattformen exisitert auch ein Druckertreiber, welcher das selbe tut. Die Erfahrung zeigt aber, dass diese teilweise unterschiedlich implementiert (z.B. unterschiedliche Dithering-Algorithmen) sind und das selbe Dokument auf verschiedenen Plattformen unterschiedlich ausgedruckt wird.

