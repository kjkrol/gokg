# Lista rzeczy do poprawy

*przekreslone zrobione

1. ~~Wydaje mi się, ze Spatial nie potrzebuje funkcji `DistanceTo(other Spatial[T], distance Distance[T]) T`~~
2. ~~Dla polygon `Bounds() Rectangle[T]` powinna byc wyliczana na rzadanie, a nie wyliczana raz w konstruktorze~~
3. ~~Funckja `Probe(margin T, plane Plane[T]) []Rectangle[T]` powinna byc czescia "systemu obliczania dystansu" BoundingBoxes.~~
4. ~~`target.Value().Probe(margin, t.plane)` jest tozsame z target.Value().Bounds().Probe(margin, t.plane); w ten sposb
   bylby tylko w Rectangle a nie w Spatial~~
5. ~~Rectanle w obenym ujeciu jest zle pojmowana struktura. Powinien sie nazywac BoundingBox i byc czescia AABB. Prawdopodobnie powinien przynalezec do gokq i razem z AABB powedrowac do gokq !~~
6. ~~Powinny pozostac tylko 3 rodzaje Spatials: Vec, Line, Polygon~~
7. ~~Usunac Line i jego transformacje.~~
8. ~~Odpuszczamy translacje shape na rzecz transformacji AABB! To AABB sa uzywane przez gokx!~~
9.  ~~Nalezy dodac mergowanie AABB !~~