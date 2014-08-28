package transparencia

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	Token string
}

const apiRoot = "http://api.transparencia.org.br/api/"

func New(token string) *Client {
	return &Client{
		Token: token,
	}
}

//func (cl *Client) Candidatos() []Candidato {
//}

//type Candidato struct {
//Id        string // (string, optional): Identificador do candidato,
//Apelido   string // (string, optional): Apelido do candidato,
//Nome      string // (string, optional): Nome do candidato,
//Numero    string // (string, optional): Número do candidato na eleição atual,
//Titulo    string // (string, optional): Número do título eleitoral do candidato,
//CPF       string // (string, optional): Número do Cadastro de Pessoas Físicas, sem formatação,
//Matricula string // (string, optional): Matrícula do candidato no sistema interno do TSE,
//Cargo     string // (string, optional): Cargo pelo qual concorre,
//Estado    string // (string, optional): Estado da federação pelo qual concorre,
//Partido   string // (string, optional): Partido pelo qual concorre,
//Idade     string // (string, optional): Idade do candidato,
//Instrucao string // (string, optional): Grau de instrução, por exemplo, 'superior completo',
//Ocupacao  string // (string, optional): Qual a ocupação atual do candidato, por exemplo, 'médico',
//MiniBio   string // (string, optional): Dados pessoais sobre o candidato (somente para candidatos a presidente, governador, deputado federal, senador pelo Paraná e por quem esteja em exercício na câmara ou senado),
//Cargos    string // (string, optional): Histórico dos cargos públicos que o candidato ocupou em sua carreira (somente para candidatos a presidente, governador, deputado federal, senador pelo Paraná e por quem esteja em exercício na câmara ou senado),
//Previsao  string // (string, optional): Custo previsto da campanha. Essa informação não é totalmente confiável,
//Bancadas  string // (string, optional): Bancadas as quais o candidato pertence. Caso não pertença a nenhuma bancada, esse atributo será vazio. Esses dados existem para candidatos a presidente, governador, deputado federal, senador pelo Paraná e por quem esteja em exercício na câmara ou senado,
//Processos string // (string, optional): Processos que o candidato responde na Justiça ou Tribunais de Contas (somente para candidatos a presidente, governador, deputado federal, senador pelo Paraná e por quem esteja em exercício na câmara ou senado). Texto formatado, com links, separados por ,
//CasaAtual string // (string, optional): Indica se o candidato está em exercício na Câmara dos Deputados ou no Senado Federal (retorna 1, se em exercício na Câmara e 2 se em exercício no Senado),
//Reeleicao bool   //(boolean, optional): Indica se o candidato está tentando a reeleição para o mesmo cargo.,
//Foto      string // (string, optional): Link para a foto oficial do candidato no site do TSE
//}

func (cl *Client) Excelencias(query map[string]string) ([]Excelencia, error) {
	var e []Excelencia
	data, err := cl.request("v1/excelencias", query)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

type Excelencia struct {
	Id                 string // (string, optional): Identificador do parlamentar (não é o mesmo id do candidato),
	Nome               string // (string, optional): Nome do parlamentar,
	Apelido            string // (string, optional): Apelido eleitoral,
	Casa               string // (string, optional): Casa a qual pertence (Câmara dos Deputados ou Senado Federal),
	Titulo             string // (string, optional): Número do título eleitoral do parlamentar,
	CPF                string // (string, optional): Número do Cadastro de Pessoas Físicas, sem formatação,
	Estado             string // (string, optional): Sigla do estado (UF) ao qual pertence,
	Partido            string // (string, optional): Sigla do partido ao qual pertence,
	MiniBio            string // (string, optional): Dados pessoais sobre o parlamentar,
	Cargos             string // (string, optional): Histórico dos cargos públicos que o parlamentar ocupou em sua carreira,
	Processos          string // (string, optional): Processos que o parlamentar responde na Justiça ou Tribunais de Contas. Texto formatado, com links, separados por ,
	Bancadas           string // (string, optional): Bancadas as quais o candidato pertence. Caso não pertença a nenhuma bancada, esse atributo será vazio.,
	PartidosPassados   string // (string, optional): Lista de partidos (siglas) pelos quais concorreu, separados por vírgula.,
	CargosPassados     string // (string, optional): Lista de cargos pelos quais concorreu, separados por vírgula.,
	EstadosPassados    string // (string, optional): Lista de estados pelos quais concorreu, separados por vírgula.,
	MunicipiosPassados string // (string, optional): Lista de municípios pelos quais concorreu (quando fizer sentido), separados por vírgula.,
	VotosPassados      string // (string, optional): Lista de quantidades de votos recebidos, separados por vírgula.,
	RecursosPassados   string // (string, optional): Lista de doações recebidas, separadas por vírgula.,
	ResultadosPassados string // (string, optional): Lista de resultados das eleições (Eleito ou Não Eleito), separados por vírgula.,
	AnosPassados       string // (string, optional): Lista dos anos das eleições que participou, separados por vírgula.
}

func (cl *Client) request(resource string, query map[string]string) ([]byte, error) {
	u, err := url.Parse(apiRoot + resource)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	log.Println(u.String())

	client := &http.Client{}
	request, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("App-Token", cl.Token)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
