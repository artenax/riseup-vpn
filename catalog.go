// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package main

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p := messageKeyToIndex[key]
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"en_US": &dictionary{index: en_USIndex, data: en_USData},
		"es":    &dictionary{index: esIndex, data: esData},
		"es_ES": &dictionary{index: es_ESIndex, data: es_ESData},
		"fr":    &dictionary{index: frIndex, data: frData},
		"nl":    &dictionary{index: nlIndex, data: nlData},
		"pt_BR": &dictionary{index: pt_BRIndex, data: pt_BRData},
	}
	fallback := language.MustParse("en-US")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"%[1]s is an easy, fast, and secure VPN service from riseup.net. %[1]s does not require a user account, keep logs, or track you in any way.\n\t    \nThis service is paid for entirely by donations from users like you. Please donate at https://riseup.net/vpn/donate.\n\t\t\nBy using this application, you agree to the Terms of Service available at https://riseup.net/tos. This service is provide as-is, without any warranty, and is intended for people who work to make the world a better place.\n\n\n%[1]v version: %[2]s": 2,
	"%s blocking internet": 26,
	"%s off":               22,
	"%s on":                21,
	"About":                3,
	"About...":             17,
	"An error has ocurred initializing %s: %v": 9,
	"Can't connect to %s: %v":                  7,
	"Cancel":                                   13,
	"Cancel connection to %s":                  14,
	"Checking status...":                       10,
	"Connecting to %s":                         23,
	"Could not find a polkit authentication agent. Please run one and try again.": 5,
	"Donate":                       1,
	"Donate...":                    16,
	"Error starting VPN":           8,
	"Help...":                      15,
	"Initialization error":         4,
	"Missing authentication agent": 6,
	"Quit":                  18,
	"Retry":                 25,
	"Route traffic through": 19,
	"Stopping %s":           24,
	"The %s service is expensive to run. Because we don't want to store personal information about you, there is no accounts or billing for this service. But if you want the service to continue, donate at least $5 each month.\n\t\nDo you want to donate now?": 0,
	"Turn off":          12,
	"Turn on":           11,
	"Use %s %v gateway": 20,
}

var en_USIndex = []uint32{ // 28 elements
	0x00000000, 0x000000fd, 0x00000104, 0x00000300,
	0x00000306, 0x0000031b, 0x00000367, 0x00000384,
	0x000003a2, 0x000003b5, 0x000003e4, 0x000003f7,
	0x000003ff, 0x00000408, 0x0000040f, 0x0000042a,
	0x00000432, 0x0000043c, 0x00000445, 0x0000044a,
	0x00000460, 0x00000478, 0x00000481, 0x0000048b,
	0x0000049f, 0x000004ae, 0x000004b4, 0x000004cc,
} // Size: 136 bytes

const en_USData string = "" + // Size: 1228 bytes
	"\x02The %[1]s service is expensive to run. Because we don't want to stor" +
	"e personal information about you, there is no accounts or billing for th" +
	"is service. But if you want the service to continue, donate at least $5 " +
	"each month.\x0a\x09\x0aDo you want to donate now?\x02Donate\x02%[1]s is " +
	"an easy, fast, and secure VPN service from riseup.net. %[1]s does not re" +
	"quire a user account, keep logs, or track you in any way.\x0a\x09    " +
	"\x0aThis service is paid for entirely by donations from users like you. " +
	"Please donate at https://riseup.net/vpn/donate.\x0a\x09\x09\x0aBy using " +
	"this application, you agree to the Terms of Service available at https:/" +
	"/riseup.net/tos. This service is provide as-is, without any warranty, an" +
	"d is intended for people who work to make the world a better place.\x0a" +
	"\x0a\x0a%[1]v version: %[2]s\x02About\x02Initialization error\x02Could n" +
	"ot find a polkit authentication agent. Please run one and try again.\x02" +
	"Missing authentication agent\x02Can't connect to %[1]s: %[2]v\x02Error s" +
	"tarting VPN\x02An error has ocurred initializing %[1]s: %[2]v\x02Checkin" +
	"g status...\x02Turn on\x02Turn off\x02Cancel\x02Cancel connection to %[1" +
	"]s\x02Help...\x02Donate...\x02About...\x02Quit\x02Route traffic through" +
	"\x02Use %[1]s %[2]v gateway\x02%[1]s on\x02%[1]s off\x02Connecting to %[" +
	"1]s\x02Stopping %[1]s\x02Retry\x02%[1]s blocking internet"

var esIndex = []uint32{ // 28 elements
	0x00000000, 0x000000ec, 0x000000f1, 0x00000305,
	0x0000030f, 0x0000030f, 0x00000376, 0x00000398,
	0x000003bb, 0x000003d3, 0x00000403, 0x0000041c,
	0x00000425, 0x0000042c, 0x00000435, 0x00000452,
	0x0000045b, 0x00000463, 0x00000470, 0x00000477,
	0x0000048a, 0x000004ab, 0x000004bb, 0x000004c9,
	0x000004dd, 0x000004f5, 0x00000500, 0x0000051a,
} // Size: 136 bytes

const esData string = "" + // Size: 1306 bytes
	"\x02El servicio %[1]s es caro de mantener. Como no queremos guardar ning" +
	"una información personal tuya, no hay cuentas ni servicio de facturación" +
	". Si quieres que este servicio continúe, dona al menos $5 cada mes.\x0a" +
	"\x09\x0a¿Quieres donar ahora?\x02Dona\x02%[1]s es un servicio de VPN fac" +
	"il, rapido y seguro de riseup.net. %[1]s no requiere registrar una cuent" +
	"a, recoge logs ni te rastrea de ninguna manera.\x0a\x09    \x0aEste serv" +
	"icio se paga completamente por donaciones de gente como tu. Por favor do" +
	"na a https://riseup.net/vpn/donate.\x0a\x09\x09\x0aAl usar este programa" +
	" estas aceptando los Terminos de servicio disponibles en https://riseup." +
	"net/tos. Este servicio se ofrece tal cual, sin garantia y con la intenci" +
	"on de la gente que trabaja en el de hacer el mundo un mejor lugar.\x0a" +
	"\x0a\x0a%[1]s version: %[2]s\x02Acerca de\x02No se pudo encontrar ningún" +
	" agente de autenticacion de polkit. Por favor lanza uno y prueba de nuev" +
	"o.\x02Falta un agente de autenticación\x02No puedo conectar con %[1]s: %" +
	"[2]v\x02Error arrancando la VPN\x02Un error ha ocurrido inicializando %[" +
	"1]s: %[2]v\x02Comprobando el estado...\x02Encender\x02Apagar\x02Cancelar" +
	"\x02Cancela la conexión a %[1]s\x02Ayuda...\x02Dona...\x02Acerca de..." +
	"\x02Cerrar\x02Salir a través de\x02Usa la salida de %[1]s por %[2]v\x02%" +
	"[1]s encendida\x02%[1]s apagada\x02Connectando a %[1]s\x02Desconnectando" +
	" de %[1]s\x02Reintentar\x02%[1]s bloqueando internet"

var es_ESIndex = []uint32{ // 28 elements
	0x00000000, 0x000000ec, 0x000000f2, 0x00000340,
	0x0000034a, 0x0000034a, 0x000003ab, 0x000003ca,
	0x000003ee, 0x00000402, 0x00000432, 0x00000448,
	0x00000450, 0x0000045b, 0x00000464, 0x0000047f,
	0x00000488, 0x00000491, 0x0000049e, 0x000004a4,
	0x000004c2, 0x000004e2, 0x000004f1, 0x00000503,
	0x00000516, 0x00000527, 0x00000532, 0x0000054c,
} // Size: 136 bytes

const es_ESData string = "" + // Size: 1356 bytes
	"\x02Correr el servicio %[1]s es caro. Porque no queremos almacenar infor" +
	"mación personal acerca tuyo, no hay cuentas o tarifas por este servicio." +
	" Pero si quieres que el mismo continúe, dona la menos USD 5 por mes.\x0a" +
	"\x09\x0a¿Quieres donar ahora?\x02Donar\x02%[1]s es un servicio de VPN fá" +
	"cil, rápido y seguro de riseup.net. %[1]s no requiere una cuenta de usua" +
	"rio, no mantiene bitácoras, o te rastrea de cualquier manera.\x0a\x09   " +
	" \x0aEl costo de este servicio está cubierto por completo por donaciones" +
	" de usuarios como tú. Por favor dona a https://riseup.net/vpn/donate." +
	"\x0a\x09\x09\x0aAl usar esta aplicación, estás de acuerdo con los Términ" +
	"os del Servicio disponibles en https://riseup.net/tos. Este servicio se " +
	"provee como está, sin ninguna garantía, y está apuntado a personas que t" +
	"rabajan para hacer del mundo un mejor lugar.\x0a\x0a\x0a%[1]v versión: %" +
	"[2]s\x02Acerca de\x02No se pudo encontrar un agente de autenticación pol" +
	"kit. Por favor corre uno e intenta de nuevo.\x02Falta agente de autentic" +
	"ación\x02No se puede conectar a %[1]s: %[2]v\x02Error iniciando VPN\x02H" +
	"a ocurrido un error inicializando %[1]s: %[2]v\x02Comprobando estado..." +
	"\x02Activar\x02Desactivar\x02Cancelar\x02Cancelar conexión a %[1]s\x02Ay" +
	"uda...\x02Donar...\x02Acerca de...\x02Salir\x02Enrutar tráfico a través " +
	"de\x02Usar ruta de salida %[1]s %[2]v\x02%[1]s activada\x02%[1]s desacti" +
	"vada\x02Conectando a %[1]s\x02Deteniendo %[1]s\x02Reintentar\x02%[1]s bl" +
	"oqueando Internet"

var frIndex = []uint32{ // 28 elements
	0x00000000, 0x00000154, 0x00000161, 0x000003d5,
	0x000003df, 0x000003f9, 0x0000045b, 0x00000487,
	0x000004b4, 0x000004d0, 0x00000513, 0x00000531,
	0x00000539, 0x00000545, 0x0000054d, 0x0000056b,
	0x00000573, 0x00000583, 0x00000590, 0x00000598,
	0x000005b0, 0x000005d4, 0x000005e6, 0x000005fc,
	0x0000060f, 0x0000061f, 0x00000628, 0x0000063e,
} // Size: 136 bytes

const frData string = "" + // Size: 1598 bytes
	"\x02L’exploitation du service %[1]s coûte cher. Dans la mesure où ne nou" +
	"s voulons enregistrer aucun renseignement personnel à votre sujet, il n’" +
	"y ni compte ni facturation pour ce service. Mais si vous souhaitez toute" +
	"fois que le service continue, faites un don d’au moins 5\u00a0$ mensuell" +
	"ement\x0a\x09\x0aSouhaitez-vous faire un don maintenant\u2009?\x02Faire " +
	"un don\x02%[1]s est un service de RPV simple, rapide et sécurisé offert " +
	"par riseup.net. %[1]s ne demande pas de compte utilisateur, ne conserve " +
	"pas de journaux, ni ne vous suit à la trace d’aucune façon.\x0a\x09    " +
	"\x0aCe service est entièrement financé par les dons d’utilisateurs comme" +
	" vous. Veuillez faire un don à https://riseup.net/vpn/donate.\x0a\x09" +
	"\x09\x0aEn utilisant cette application, vous acceptez les conditions gén" +
	"érales d’utilisation qui se trouvent sur https://riseup.net/tos. Ce ser" +
	"vice est fourni tel quel, sans aucune garantie, et s’adresse aux personn" +
	"es qui œuvrent à rendre le monde meilleur.\x0a\x0a\x0a%[1]v version " +
	"\u00a0: %[2]s\x02À propos\x02Erreur d’initialisation\x02Impossible de tr" +
	"ouver un agent d’authentification polkit. Veuillez en exécuter un et res" +
	"sayer.\x02L’agent d’authentification est manquant\x02Impossible de se co" +
	"nnecter à %[1]s\u00a0: %[2]v\x02Erreur du démarrage du RPV\x02Une erreur" +
	" est survenue lors de l'initialisation de %[1]s\u00a0: %[2]v\x02Vérifica" +
	"tion de l’état…\x02Activer\x02Désactiver\x02Annuler\x02Annuler la connex" +
	"ion à %[1]s\x02Aide…\x02Faire un don…\x02À propos…\x02Quitter\x02Achemin" +
	"er le trafic par\x02Utiliser la passerelle %[1]s %[2]v\x0a\x02%[1]s est " +
	"activé\x02%[1]s est désactivé\x02Connexion à %[1]s\x02Arrêt de %[1]s\x02" +
	"Ressayer\x02%[1]s bloque Internet"

var nlIndex = []uint32{ // 28 elements
	0x00000000, 0x0000010c, 0x00000114, 0x00000356,
	0x0000035b, 0x0000036d, 0x000003bc, 0x000003d9,
	0x000003fd, 0x0000041f, 0x00000468, 0x0000047e,
	0x0000048a, 0x00000497, 0x000004a1, 0x000004bf,
	0x000004c7, 0x000004d2, 0x000004da, 0x000004e2,
	0x000004f5, 0x00000511, 0x0000051b, 0x00000525,
	0x00000539, 0x0000054f, 0x00000560, 0x0000057d,
} // Size: 136 bytes

const nlData string = "" + // Size: 1405 bytes
	"\x02De %[1]s dienst is kostbaar om te onderhouden. Omdat we geen persoon" +
	"lijke informatie over u willen bijhouden, zijn er geen accounts of betal" +
	"ingen voor deze dienst. Om deze dienst in leven te houden, overweeg ten " +
	"minste €4 per maand te schenken.\x0a\x09\x0aWilt u nu doneren?\x02Donere" +
	"n\x02%[1]s is een gemakkelijke, snelle, en veilige VPN-dienst van riseup" +
	".net. %[1]s vereist geen gebruikersaccount, houdt geen logboek bij en vo" +
	"lgt u niet op wat voor manier dan ook.\x0a\x09\x0aDeze dienst wordt voll" +
	"edig gefinancierd door donaties van gebruikers zoals u. Overweeg bij te " +
	"dragen op https://riseup.net/vpn/donate.\x0a\x09\x09\x0aDoor deze applic" +
	"atie te gebruiken gaat u akkoord met onze gebruikersvoorwaarden die besc" +
	"hikbaar zijn op https://riseup.net/tos. Deze dienst wordt geleverd zonde" +
	"r enige garantie, en is bedoeld voor mensen die werken aan een betere we" +
	"reld.\x0a\x0a\x0a%[1]v versie: %[2]s\x02Over\x02Initialisatiefout\x02Kan" +
	" geen polkit authenticatieagent vinden. Voer er een uit en probeer opnie" +
	"uw.\x02Authenticatieagent ontbreekt\x02Kan niet verbinden met %[1]s: %[2" +
	"]v\x02Fout bij het opstarten van de VPN\x02Er heeft zich een fout voorge" +
	"daan bij het initialiseren van %[1]s: %[2]v\x02Status controleren...\x02" +
	"Inschakelen\x02Uitschakelen\x02Annuleren\x02Annuleer verbinding met %[1]" +
	"s\x02Hulp...\x02Doneren...\x02Over...\x02Stoppen\x02Stuur verkeer door" +
	"\x02Gebruik %[1]s %[2]v gateway\x02%[1]s aan\x02%[1]s uit\x02Verbinden m" +
	"et %[1]s\x02%[1]s aan het stoppen\x02Opnieuw proberen\x02%[1]s blokkeert" +
	" het internet"

var pt_BRIndex = []uint32{ // 28 elements
	0x00000000, 0x000000fd, 0x00000110, 0x0000030c,
	0x00000312, 0x0000032a, 0x0000039b, 0x000003c6,
	0x000003e4, 0x000003fa, 0x00000429, 0x0000043f,
	0x00000445, 0x0000044e, 0x00000457, 0x00000474,
	0x0000047d, 0x00000493, 0x0000049c, 0x000004a1,
	0x000004bf, 0x000004dd, 0x000004f0, 0x00000506,
	0x00000523, 0x0000053e, 0x0000054f, 0x00000571,
} // Size: 136 bytes

const pt_BRData string = "" + // Size: 1393 bytes
	"\x02The %[1]s service is expensive to run. Because we don't want to stor" +
	"e personal information about you, there is no accounts or billing for th" +
	"is service. But if you want the service to continue, donate at least $5 " +
	"each month.\x0a\x09\x0aDo you want to donate now?\x02Fazer uma doação" +
	"\x02%[1]s is an easy, fast, and secure VPN service from riseup.net. %[1]" +
	"s does not require a user account, keep logs, or track you in any way." +
	"\x0a\x09    \x0aThis service is paid for entirely by donations from user" +
	"s like you. Please donate at https://riseup.net/vpn/donate.\x0a\x09\x09" +
	"\x0aBy using this application, you agree to the Terms of Service availab" +
	"le at https://riseup.net/tos. This service is provide as-is, without any" +
	" warranty, and is intended for people who work to make the world a bette" +
	"r place.\x0a\x0a\x0a%[1]v version: %[2]s\x02Sobre\x02Erro na inicializaç" +
	"ão\x02Não foi possível encontrar um agente de autenticação polkit. Por " +
	"favor, execute um agente e tente novamente.\x02Um agente de autenticação" +
	" está faltando\x02Can't connect to %[1]s: %[2]v\x02Erro ao iniciar a VPN" +
	"\x02An error has ocurred initializing %[1]s: %[2]v\x02Verificando estado" +
	"...\x02Ligar\x02Desligar\x02Cancelar\x02Cancelar a conexão à %[1]s\x02Aj" +
	"uda...\x02Fazer uma doação...\x02Sobre...\x02Sair\x02Rotear o tráfego at" +
	"ravés de\x02Usar o gateway %[2]v da %[1]s\x02%[1]s está ligada\x02%[1]s " +
	"está desligada\x02A %[1]s está sendo iniciada\x02A %[1]s está sendo para" +
	"da\x02Tentar novamente\x02%[1]s está bloqueando a Internet"

	// Total table size 9102 bytes (8KiB); checksum: 2ACD6A17
