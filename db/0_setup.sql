CREATE TABLE "events" (
    "id" SERIAL PRIMARY KEY,
    "public_id" varchar(12),
    CONSTRAINT "idx_public_id" UNIQUE ("public_id"),
    "name" varchar(30) NOT NULL,
    "sport" varchar(20) NOT NULL,
    "date" DATE,
    "location" varchar(50) NOT NULL,
    "price" smallint NOT NULL DEFAULT 0,
    "description" TEXT NOT NULL,
    "level" varchar(30) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "owner_id" INT NOT NULL
);

CREATE TABLE "event_levels" ("id" int PRIMARY KEY, "name" varchar(20));
INSERT INTO "event_levels" ("id", "name") VALUES (1, 'zacatecnik'), (2, 'pokrocily'), (3, 'expert');

CREATE TABLE "sports" ("id" int PRIMARY KEY, "name" varchar(20));
INSERT INTO "sports" ("id", "name") VALUES (1, 'basketbal'), (2, 'florbal'), (3, 'fotbal');


CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(30) NOT NULL,
  "email" varchar(30) NOT NULL,
  "password" varchar(60) NOT NULL,
  "rating" int NOT NULL
);

INSERT INTO events (id, name, sport, date, location, price, description, level, public_id, created_at, owner_id) 
VALUES 
(31, 'Pražská Cyklistická Grand Prix', 'Cyklistika', '2024-04-01', 'Praha', 450, 'Pražská Cyklistická Grand Prix" je vzrušující mezinárodní cyklistický závod, který se koná v srdci České republiky, v historické Praze. Tento jedinečný sportovní event, rozložený na několik dní, zahrnuje několik etap, které se rozprostírají přes různé části města a jeho okolí, a nabízí jedinečný mix městského a přírodního terénu. <br><br> Závodníci se vydají na cestu přes známé pražské památky, jako jsou Karlův most, Pražský hrad a Staroměstské náměstí, čímž získají nevšední pohled na město z pohledu cyklistiky. Okolo trati se budou střídat historické uličky, parky a říční břehy Vltavy, což představuje výzvy v podobě různorodých povrchů a terénů. <br><br> Kromě sportovního zápolení, "Pražská Cyklistická Grand Prix" nabízí bohatý doprovodný program, včetně kulturních akcí, koncertů a výstav, které oslavují českou kulturu a historii. Tyto akce poskytují účastníkům i divákům možnost prohloubit své poznání českého dědictví a užít si místní gastronomii a pohostinnost. <br><br> Závod také zdůrazňuje důležitost udržitelnosti a ekologie, propaguje používání jízdních kol jako ekologického způsobu dopravy ve městech a podporuje různé lokální environmentální iniciativy. <br><br> "Pražská Cyklistická Grand Prix" je tedy nejen sportovní událostí, ale i oslavou Prahy, její kultury a historie, a přitahuje jak profesionální cyklisty, tak cyklistické nadšence a turisty z celého světa, kteří chtějí prožít jedinečnou atmosféru tohoto historického města.', 'Expert', 'hsn45lkm7te8', '2023-11-03 11:52:14', 0),
(46, 'Basketball Match at Park', 'Basketball', '2023-12-29', 'Central Park', 123, 'Připojte se k nám na vzrušující basketbalový turnaj, který se koná na nově zrekonstruovaném venkovním hřišti. Tento event je určen pro hráče všech věkových kategorií a úrovní dovedností. Čekají na vás zápasy 3 na 3, dovednostní soutěže a skvělá atmosféra.', 'Advanced', '6li0xwokov08', '2023-11-04 19:07:59', 0),
(76, 'Harmonie Pohybu: Yoga Retreat', 'Yoga', '2024-02-28', 'Olomouc', 99, 'Harmonie v Pohybu je mezinárodní yoga retreat, který se koná v malebné přírodní lokalitě. Tento víkendový event nabízí účastníkům možnost prohloubit svou praxi yogy, meditace a mindfulness pod vedením světově uznávaných instruktorů. Program zahrnuje různé styly yogy, od dynamické vinyasy po uklidňující yin yogu, doplněný o workshopy na téma duševního zdraví a wellness. Účastníci mají příležitost se spojit s přírodou, obnovit svou energii a naučit se techniky pro zlepšení svého každodenního života.', 'Any', '6kqoucaz8jcy', '2023-11-27 11:48:27', 0);

INSERT INTO events (id, name, sport, date, location, price, description, level, public_id, created_at, owner_id) 
VALUES 
(77, 'Rytmus Metropole, Dance Fest', 'Tanec', '2024-04-30', 'Ostrava', 249, 'Rytmus Metropole je dynamický městský tančírní festival, který se koná na různých venkovních a indoorových místech ve městě. Festival představuje širokou škálu tanečních stylů, od street dance až po salsa, bachata, a moderní jazz. Tento víkendový festival nabízí workshopy pro tanečníky všech úrovní, vedené profesionálními tanečníky a choreografy. Večery jsou plné živých vystoupení, tanečních soutěží a společenských tanců, což dává účastníkům příležitost ukázat své dovednosti a zapojit se do taneční komunity.', 'Advanced', 'kxov5znwscvq', '2023-11-27 11:52:06', 0),
(78, 'Síla a Stamina', 'Fitness', '2024-06-25', 'Pardubice', 349, 'Síla a Stamina" je víkendová fitness challenge, která se zaměřuje na posilování a vytrvalost. Tato událost zahrnuje řadu soutěží a výzev, jako jsou váhové kategorie v deadlift, squat a bench press, spolu s kardiovými výzvami jako sprinty a překážkové běhy. Event je otevřený pro účastníky všech fitness úrovní, od amatérů po profesionály, a nabízí personalizované tréninkové plány a semináře o výživě a regeneraci. Je to skvělá příležitost pro fitness nadšence setkat se, soutěžit a naučit se nové techniky pro zlepšení fyzické kondice.', 'Expert', 'nd86fdmgkbc9', '2023-11-27 11:53:53', 0),
(79, 'Bojovníkův Summit', 'Ju Jitsu', '2024-06-25', 'Pardubice', 500, 'Bojovníkův Summit" je prestižní mezinárodní turnaj v ju jitsu, který přitahuje bojovníky z celého světa. Tento turnaj nabízí kategorie pro různé pásy a věkové skupiny, od juniorů po seniory, a zahrnuje jak tradiční, tak moderní ju jitsu. Kromě soutěžních zápasů, event zahrnuje semináře a workshopy vedené veterány a expertními trenéry, kde se účastníci mohou naučit pokročilé techniky a strategie. Turnaj je také oslavou bojového umění a jeho filozofie, zdůrazňující důležitost disciplíny, respektu a neustálého zlepšování.', 'Expert', 'e5dozys0skhv', '2023-11-27 11:55:18', 0);


INSERT INTO events (id, name, sport, date, location, price, description, level, public_id, created_at, owner_id) 
VALUES 
(81, 'Zelená Výzva', 'Golf', '2023-11-29', 'Praha', 568, '"Zelená Výzva: Pražský Golf Open" je prestižní golfový turnaj, který se koná v okolí Prahy, srdci České republiky. Tento turnaj přiláká golfisty různých úrovní, od amatérů až po profesionály, kteří se utkají na několika špičkových golfových hřištích v regionu. Hřiště jsou známá svými náročnými fairwayi a rafinovaně navrženými greeny, což slibuje vzrušující a strategickou hru. Kromě samotné soutěže se konají i různé doprovodné akce, včetně tréninkových seminářů s profesionálními trenéry, networkingových večerů a galavečeří. "Zelená Výzva" je také příležitostí pro účastníky a návštěvníky, aby si užili krásu a historii Prahy, a zároveň se zapojili do nadšené golfové komunity.', 'Expert', 'tvlitx9jb57v', '2023-11-27 18:09:25', 0),
(82, 'Brněnský Sportovní Běh', 'Běh', '2023-11-30', 'Brno', 344, 'Přijďte s přáteli a zažijte vzrušující den plný sportu a dobrodružství v rámci <b>Brněnského Sportovního Běhu! </b>Tato událost je otevřena všem věkovým kategoriím, takže vás zveme, abyste se připojili a vytvořili týmy pro štafetový závod. <br><br> Prozkoumáte jeden z našich krásných brněnských parků, závodíte a soupeříte s ostatními týmy, ale také se těšíte na společný čas strávený v přírodě. <br><br> Spolupráce a týmový duch jsou klíčovými slovy této události, takže přijďte a podpořte se navzájem v dosažení společného cíle.', 'Any', 'ymy16bg7ozh3', '2023-11-27 19:08:55', 0),
(85, 'Name example', 'Sport example', '2023-11-03', 'Definitely not Brno', 1, 'Description example', 'Any', 'mrgldgdohxpe', '2024-03-06 15:51:50', 51);
