[] - notify
[X] - Get what i want from the json
[X] - Send a email
[] - Create a database Schema
[] - Create User
[] - What to notify
[X] - A way to choose what i want the repo notify
[X] - Create a webhooks using github
[X] - Using a exist github token KEY
[X] - Creating a new github token KEY
[X] - Create a webhooks for others users
[X] - Test with a real webhook

criar um banco de dados, motivos:
tabela instalation
tabela user
tabela repositorios
tabela eventos
tabela tipo de eventos?
preciso colocar o instalation ID e associar com usuario
preciso adicionar os usuarios
preciso saber quais repositorios do usuario ele esta usando meu app
crio um JWT salvo no bd pra saber se ainda esta ativo?

rota de webhooks:
[X] - rota para receber cada requisição da instalacao do APP - aproveitar e salvar no banco de dados todos os repositorios que ele liberou para mim

rotas no geral:
OAUTH E TA BOM? CRUD?: usuario
CRUD: webhook - Criar um webhook - Editar um webhook: - eventos - Deletar um webhook

TOKEN JWT:
nao preciso criar um sempre para cada nova pessoa ja que tem a ver com a minha chave e meu APP ID, criar uma forma de validacao pra ve se ele ainda esta valido e se passar do tempo eu crio um novo
