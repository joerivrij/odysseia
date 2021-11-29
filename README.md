# Odysseia <!-- omit in toc -->

Odysseia (Ὀδύσσεια) is one of the two famous poems bij Homeros. It describes the journey of Odysseus and his crew to get home. Learning Greek is a bit like that - a odyssey. It is a hobby project that combines a few of my passions, both ancient Greek (history) and finding technical solutions for problems. As This is a hobby project first and foremost any mistakes are my own, either in translation or in interpation of text.

The goal is for people to learn or rehearse ancient Greek. Some of it is in Dutch but most of it is in English. There is also a dictionary that you can search through. Most of it is still very much a work in progress.

# Table of contents <!-- omit in toc -->
- [Backend](#backend)
  - [Alexandros - Αλέξανδρος](#alexandros---αλέξανδρος)
  - [Dionysos - Διονύσιος ὁ Θρᾷξ](#dionysos---διονύσιος-ὁ-θρᾷξ)
  - [Herodotos - Ἡρόδοτος](#herodotos---ἡρόδοτος)
  - [Sokrates - Σωκράτης](#sokrates---σωκράτης)
  - [Solon - Σόλων](#solon---σόλων)
- [Common](#common)
  - [Eratosthenes - Ἐρατοσθένης](#eratosthenes---ἐρατοσθένης)
  - [Plato - Πλάτων](#plato---πλάτων)
- [Dataseeders](#dataseeders)
  - [Anaximander - Ἀναξίμανδρος](#anaximander---ἀναξίμανδρος)
  - [Demokritos - Δημόκριτος](#demokritos---δημόκριτος)
  - [Herakleitos - Ἡράκλειτος](#herakleitos---ἡράκλειτος)
  - [Parmenides - Παρμενίδης](#parmenides---παρμενίδης)
- [Docs](#docs)
  - [Ploutarchos - Πλούταρχος](#ploutarchos---πλούταρχος)
- [Frontend](#frontend)
  - [Pheidias - Φειδίας](#pheidias---φειδίας)
- [Init](#init)
  - [Periandros - Περίανδρος](#periandros---περίανδρος)
- [Ops](#ops)
  - [Archimedes - Ἀρχιμήδης](#archimedes---ἀρχιμήδης)
  - [Lykourgos - Λυκοῦργος](#lykourgos---λυκοῦργος)
  - [Themistokles - Θεμιστοκλῆς](#themistokles---θεμιστοκλῆς)
- [Sidecar](#sidecar)
  - [Ptolemaios - Πτολεμαῖος](#ptolemaios---πτολεμαῖος)
- [Tests](#tests)
  - [Hippokrates - Ἱπποκράτης](#hippokrates---ἱπποκράτης)
  - [Xerxes - Ξέρξης](#xerxes---ξέρξης)

## Backend

### Alexandros - Αλέξανδρος

<img src="https://upload.wikimedia.org/wikipedia/commons/5/59/Alexander_and_Bucephalus_-_Battle_of_Issus_mosaic_-_Museo_Archeologico_Nazionale_-_Naples_BW.jpg" alt="Alexandros" width="200"/>

What could I ever say in a few lines that would do justice to one of the most influential people of all time? Alexandros's energy and search for the end of the world was relentless, so too is his search for Greek words within odysseia.

### Dionysos - Διονύσιος ὁ Θρᾷξ

<img src="https://alchetron.com/cdn/dionysius-thrax-73e8d598-e6d3-4f5f-bb04-debff25a456-resize-750.jpeg" alt="Dionysos" width="200"/>

Probably the first Greek Grammarian who wrote the "Τέχνη Γραμματική". Even though often called "the Thracian" he was most likely from Alexandria which was the hub for Greek learning for a long time.

### Herodotos - Ἡρόδοτος

Ἡροδότου Ἁλικαρνησσέος ἱστορίης ἀπόδεξις ἥδε - This is the display of the inquiry of Herodotos of Halikarnassos

<img src="https://upload.wikimedia.org/wikipedia/commons/6/6f/Marble_bust_of_Herodotos_MET_DT11742.jpg" alt="Sokrates" width="200"/>

Herodotos is often hailed as the father of history. I name he lives up to. His work (the histories) is a lively account of the histories of the Greeks and Persians and how they came into conflict. This API is responsible for passing along sentences you need to translate. They are then checked for accuracy.

### Sokrates - Σωκράτης

ἓν οἶδα ὅτι οὐδὲν οἶδα - I know one thing, that I know nothing

<img src="https://upload.wikimedia.org/wikipedia/commons/2/25/Raffael_069.jpg" alt="Sokrates" width="200"/>

Sokrates (on the right) is a figure of mythical propertions. He could stare at the sky for days, weather cold in nothing but a simple cloak. Truly one of the greatest philosophers and a big influence on Plato which is why we know so much about him at all. A sokratic dialogue is a game of wits were the back and forth between Sokrates and whoever was unlucky (or lucky) to be part of the dialogue would end in frustration. Sokrates was known to question anyone until he had proven they truly knew nothing. As the API responsible for creating and asking questions he was the obvious choice.

### Solon - Σόλων

αὐτοὶ γὰρ οὐκ οἷοί τε ἦσαν αὐτὸ ποιῆσαι Ἀθηναῖοι: ὁρκίοισι γὰρ μεγάλοισι κατείχοντο δέκα ἔτεα χρήσεσθαι νόμοισι τοὺς ἄν σφι Σόλων θῆται - since the Athenians themselves could not do that, for they were bound by solemn oaths to abide for ten years by whatever laws Solon should make

<img src="https://upload.wikimedia.org/wikipedia/commons/1/12/Ignoto%2C_c.d._solone%2C_replica_del_90_dc_ca_da_orig._greco_del_110_ac._ca%2C_6143.JPG" alt="Solon" width="200"/>

Solon is most famous for his role as the great Athenian lawgiver following the reforms made by Drakon. His laws laid the foundation of what would become the Athenian Democracy.

## Common

### Eratosthenes - Ἐρατοσθένης

<img src="https://upload.wikimedia.org/wikipedia/commons/b/b3/Eratosthene.01.png" alt="Eratosthenes" width="200"/>

Holds fixtures for testing. Eratosthenes was one of the librarians of Alexandria. He is most famous for calculating the circumference of the earth.

### Plato - Πλάτων

χαλεπὰ τὰ καλά - good things are difficult to attain

<img src="https://upload.wikimedia.org/wikipedia/commons/4/4a/Platon.png" alt="Plato" width="200"/>

## Dataseeders

### Anaximander - Ἀναξίμανδρος

οὐ γὰρ ἐν τοῖς αὐτοῖς ἐκεῖνος ἰχθῦς καὶ ἀνθρώπους, ἀλλ' ἐν ἰχθύσιν ἐγγενέσθαι τὸ πρῶτον ἀνθρώπους ἀποφαίνεται καὶ τραφέντας, ὥσπερ οἱ γαλεοί, καὶ γενομένους ἱκανους ἑαυτοῖς βοηθεῖν ἐκβῆναι τηνικαῦτα καὶ γῆς λαβέσθαι.

He declares that at first human beings arose in the inside of fishes, and after having been reared like sharks, and become capable of protecting themselves, they were finally cast ashore and took to land

<img src="https://upload.wikimedia.org/wikipedia/commons/3/38/Anaximander.jpg" alt="Anaximander" width="200"/>

Anaximander developed a rudimentary evolutionary explanation for biodiversity in which constant universal powers affected the lives of animals

### Demokritos - Δημόκριτος

νόμωι (γάρ φησι) γλυκὺ καὶ νόμωι πικρόν, νόμωι θερμόν, νόμωι ψυχρόν, νόμωι χροιή, ἐτεῆι δὲ ἄτομα καὶ κενόν 

By convention sweet is sweet, bitter is bitter, hot is hot, cold is cold, color is color; but in truth there are only atoms and the void.

<img src="https://upload.wikimedia.org/wikipedia/commons/5/58/Rembrandt_laughing_1628.jpg" alt="Demokritos" width="200"/>

Most famous for his theory on atoms, everything can be broken down into smaller parts.

### Herakleitos - Ἡράκλειτος

πάντα ῥεῖ - everything flows

<img src="https://upload.wikimedia.org/wikipedia/commons/6/67/Raphael_School_of_Athens_Michelangelo.jpg" alt="Parmenides" width="200"/>

Herakleitos is one of the so-called pre-socratics. Philosophers that laid the foundation for the future generations. One of his most famous sayings is "No man ever steps in the same river twice". Meaning everything constantly changes. Compare that to Parmenides. He is said to be a somber man, perhaps best reflected in the School of Athens painting where his likeness is taken from non other than Michelangelo.

### Parmenides - Παρμενίδης

τό γάρ αυτο νοειν έστιν τε καί ειναι - for it is the same thinking and being

<img src="https://upload.wikimedia.org/wikipedia/commons/2/20/Sanzio_01_Parmenides.jpg" alt="Parmenides" width="200"/>

Parmenides is one of the so-called pre-socratics. Philosophers that laid the foundation for the future generations. One of the key elements in his work is the fact that everything is one never changing thing. Therefor he is a good fit for the dataseeder. Making it like nothing every changed.

## Docs

### Ploutarchos - Πλούταρχος
<img src="https://upload.wikimedia.org/wikipedia/commons/0/02/Plutarch_of_Chaeronea-03-removebg-preview.png" alt="Ploutarchos" width="400"/>

Ploutarchos (or Plutarch) is most well known for his Parallel Lives, a series of books where he compares a well known Roman to a Greek counterpart.

## Frontend

### Pheidias - Φειδίας

<img src="https://upload.wikimedia.org/wikipedia/commons/d/d7/Charles_B%C3%A9ranger_-_Replica_of_The_H%C3%A9micycle_-_Walters_3783.jpg" alt="Pheidias" width="400"/>

Pheidias (or Phidias) is one of the great artists of the Greek world, most famous for his work on the Athenian Akropolis. An apt choice for the frontend of the app.

## Init

### Periandros - Περίανδρος

Περίανδρος δὲ ἦν Κυψέλου παῖς οὗτος ὁ τῷ Θρασυβούλῳ τὸ χρηστήριον μηνύσας· ἐτυράννευε δὲ ὁ Περίανδρος Κορίνθου - Periander, who disclosed the oracle's answer to Thrasybulus, was the son of Cypselus, and sovereign of Corinth

<img src="https://upload.wikimedia.org/wikipedia/commons/4/48/Periander_Pio-Clementino_Inv276.jpg" alt="Periandros" width="200"/>

Tyrant of Corinth.

## Ops

### Archimedes - Ἀρχιμήδης

εὕρηκα - I found it!

<img src="https://upload.wikimedia.org/wikipedia/commons/c/c5/Archimedes_The_School_of_Athens.png" alt="Archimedes" width="200"/>

Archimedes is one of the greatest mathematicians of all time. He is also known for some nifty inventions which is why his
name has been chosen for the `ctl` tooling.

### Lykourgos - Λυκοῦργος

<img src="https://upload.wikimedia.org/wikipedia/commons/5/57/Lycurgus.jpg" alt="Lykourgos" width="200"/>

Lykourgos (or Lycurgus) is a semi mythical lawmaker that laid the foundation for the Spartan society with strict rules. As a deployment abides by strict rules this is a great fit.

### Themistokles - Θεμιστοκλῆς

<img src="https://upload.wikimedia.org/wikipedia/commons/8/86/Illustrerad_Verldshistoria_band_I_Ill_116.png" alt="Themistokles" width="200"/>

Themistokles is argueably the greatest Greek admiral. His victory at Salamis is most well-known. As an admiral he held sway over many ships and thus over many pilots (kubernetes).


## Sidecar

### Ptolemaios - Πτολεμαῖος

<img src="https://upload.wikimedia.org/wikipedia/commons/2/21/Ptolemy_I_Soter_Louvre_Ma849.jpg" alt="Ptolemaios" width="200"/>


First Macedonian king of Egypt.

## Tests

### Hippokrates - Ἱπποκράτης

<img src="https://upload.wikimedia.org/wikipedia/commons/7/7c/Hippocrates.jpg" alt="Hippokrates" width="200"/>


The most well known medical professional in history. Hippokrates houses tests to see whether the other services are in good health.

### Xerxes - Ξέρξης

<img src="https://upload.wikimedia.org/wikipedia/commons/6/64/National_Museum_of_Iran_Darafsh_%28785%29.JPG" alt="Xerxes" width="200"/>

Xerxes was the Persian great king and invaded Greece during the second Greco-Persion war. He tested the Greeks with an army so great it was reported to drain rivers and shake the earth.