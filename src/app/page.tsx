'use client'
import { Text, Heading, Divider, Container, Box, Flex  } from '@chakra-ui/react';

export default function Home() {
  return (
    <Box w='100%' height='100%' p={10}>
      <Container p={1} centerContent>
        <Heading size='3xl' as='h1' noOfLines={1}>Bem Vindo</Heading>
        <Divider />
        <Text margin={1}>
          Bem vindo,
          para iniciar seus experimentos, instancie os workers e registre o orquestrator no menu ao lado
        </Text>
      </Container>
      <Container p={1} centerContent>
        <Heading>A Aplicação</Heading>
        <Divider />
        <Flex textAlign='justify' flexDirection='column' padding={1}>
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
        </Flex>
      </Container>
    </Box>  
  );
};
