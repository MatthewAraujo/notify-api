rota de webhooks:
[] - Usuario pode adicionar novos repositorios do github
json escrito: "added"

Dentro da rota de installation verificar o payload para saber se é "created" ou "added"
quando for created cria um installation ID
quando for added adiciona novo repositorios
[X] - quando for deleted apagar tudo deste usuario - checar fluxo inteiro porem acho que ta funfando

refazer a rota de notification pois nao fiz nenhuma logica incluindo o github, apenas no banco de dados
rota de edit
[x] - rota de delete

criar payload para deletar um webhook
criar payload para editar um webhook

USER
Fazer logica de deletar usuario
Fazer a logica de criar usuario com OAUTH - ver video do mewjke https://youtu.be/iHFQyd__2A0?si=pW82GF-D9WYu34vK

para o front
rota para pegar todos o eventos
rota para pegar todos os repositorios

regras de negocio
se um usuario for deletado, automaticamente devera revogar o acesso a conta dele portanto a tabela de installation tera o revoged_at e os webhooks seram apagados
[] - Usuario pode adicionar novos repositorios do github
