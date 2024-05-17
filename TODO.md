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

banco de dados:
[]-revisar schema

[]-crio um JWT salvo no bd pra saber se ainda esta ativo?
Criar no banco de dados uma tabela para o JWT
rota de webhooks:

[]-salvar o access token do usuario no banco para que eu nao precise
criar varios tokens
- ver se o acess token fica diferente a cada repositorio novo que o usuario
me da acesso

rotas no geral:
OAUTH E TA BOM? CRUD?: usuario
CRUD: webhook - Criar um webhook - Editar um webhook: - eventos - Deletar um webhook
[] - Usuario pode adicionar novos repositorios do github

TOKEN JWT:
nao preciso criar um sempre para cada nova pessoa ja que tem a ver com a minha chave e meu APP ID, criar uma forma de validacao pra ve se ele ainda esta valido e se passar do tempo eu crio um novo

para o front
rota para pegar todos o eventos
rota para pegar todos os repositorios

regras de negocio
se um usuario for deletado, automaticamente devera revogar o acesso a conta dele portanto a tabela de installation tera o revoged_at e os webhooks seram apagados
[] - Usuario pode adicionar novos repositorios do github
