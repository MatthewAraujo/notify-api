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

jwt criando varios

request: bad credentials resolve quando chegar em casa

Dentro da rota de installation verificar o payload para saber se é "created" ou "added"
quando for created cria um installation ID
quando for added adiciona novo repositorios

AUTH
-ver se o acess token fica diferente a cada repositorio novo que o usuario
me da acesso

rota de webhooks:
CRUD: webhook - Criar um webhook - Editar um webhook: - eventos - Deletar um webhook
[] - Usuario pode adicionar novos repositorios do github
json escrito: "added"

USER
Fazer logica de deletar usuario
Fazer a logica de criar usuario com OAUTH - ver video do mewjke https://youtu.be/iHFQyd__2A0?si=pW82GF-D9WYu34vK

para o front
rota para pegar todos o eventos
rota para pegar todos os repositorios

regras de negocio
se um usuario for deletado, automaticamente devera revogar o acesso a conta dele portanto a tabela de installation tera o revoged_at e os webhooks seram apagados
[] - Usuario pode adicionar novos repositorios do github
