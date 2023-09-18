'use client'
import { Text, Container } from '@chakra-ui/react';
import Article from './components/article';

export default function Home() {
  return (
    <>
      <Container p={1} centerContent>
        <Article title={ {text: 'Bem Vindo', as: 'h1', size: '3xl'} }>
          <Text margin={1}>
            Bem vindo,
            para iniciar seus experimentos, instancie os workers e registre o orquestrator no menu ao lado
          </Text>
        </Article>
      </Container>
      <Container p={1} centerContent>
        <Article title={ {text: 'A Aplicação'} }>
            <Text margin={1}>
              O MQTTPoller é uma interface para o MQTTLoader, que tem como objetivo facilitar o disparo de experimentos e a coleta de resultados. Para isso, utiliza uma arquitetura de orquestrador
              e trabalhadores que controlam a execução dos experimentos. Para isso, são utilizados os protocolos HTTP e MQTT, um para a comunicação com o usuário e outro para comunicação interna. A 
              aplicação permite o disparo de multiplos experimentos em paralelo, seleção dos trabalhadores que executaram os experimentos e o cancelamento de experimentos em execução.
            </Text>
            <Text margin={1}>
              Os orquestradores, são os nodos coordenadores da aplicação. São responsáveis por: servir de interface com o usuário, registrar os trabalhadores, controlar quais experimentos serão
              executados em cada trabalhador e armazenar o registro de quais experimentos já foram executados e seus resultados. Para a comunicação com os trabalhadores, se utiliza o protocolo
              MQTT, baseado em publicação e assinatura. Para a interface com o usuário, é utilizado uma API Rest com o protocolo HTTP, com as rotas definidas na documentação do orquestrador.
            </Text>
            <Text margin={1}>
              Por fim, os trabalhadores são os responsáveis por executar os experimentos e coletar os seus resultados. Para isso, eles executam o MQTTLoader com os parâmetros passados pelo orquestrador.
              Ao finalizar a execução do experimentos, os seus resultados são enviados via MQTT no formato JSON, com a identificação de qual experimentos eles estavam executando e as métricas do
              MQTTLoader. São totalmente desacloplados entre si, entretando, dependentes do orquestrador.
            </Text>
        </Article>
      </Container>
    </>  
  );
};
