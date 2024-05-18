checar duplicidade de tudo!!

rota de webhooks:
CRUD: webhook - Criar um webhook - Editar um webhook: - eventos - Deletar um webhook
[] - Usuario pode adicionar novos repositorios do github
json escrito: "added"

Dentro da rota de installation verificar o payload para saber se é "created" ou "added"
quando for created cria um installation ID
quando for added adiciona novo repositorios

USER
Fazer logica de deletar usuario
Fazer a logica de criar usuario com OAUTH - ver video do mewjke https://youtu.be/iHFQyd__2A0?si=pW82GF-D9WYu34vK

para o front
rota para pegar todos o eventos
rota para pegar todos os repositorios

regras de negocio
se um usuario for deletado, automaticamente devera revogar o acesso a conta dele portanto a tabela de installation tera o revoged_at e os webhooks seram apagados
[] - Usuario pode adicionar novos repositorios do github
