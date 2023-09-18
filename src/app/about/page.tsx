'use client'

import { Container, Divider, Heading, Highlight, Text, Flex } from "@chakra-ui/react";
import Article from "../components/article";

const Homepage = () => {
    return (
        <>
            <Container>
                <Article title={ { text: 'Sobre o Projeto' } }>
                    <Text margin={1}>
                        Esse projeto surgiu a partir de um projeto de iniciação científica chamado 
                        Avaliação da escalabilidade do protocolo MQTT no contexto da Internet das Coisas, financiado pela UFFS,
                        com o código PES-2021-0471. A ideia inicial, era estender uma ferramenta de testes para redes MQTT, mas o rumo dele mudou
                        para melhorar a usabilidade do MQTTLoader.
                    </Text>
                    <Text margin={1}>
                        Um grande problema das ferramentas estudas era a dificuldade de gerar experimentos com clientes hospedados em máquinas físicas diferentes.
                        No geral, a maioria se preocupa em gerar um trafégo de rede local, o que limita a quantidade de testes que podem ser gerados com elas. As
                        ferramentas que não possuiam essa limitação eram o MQTTLoader e o plugin para o JMQTER, MQTT-JMETER. Por fim, foi escolhido melhorar a
                        usabilidade do MQTTLoader, com o objetivo de atrair mais pessoas para a comunidade.
                    </Text>
                    <Text margin={1}>
                        O grande problema do MQTTLoader, é a falta de uma interface gráfica, por ser um programa de linha de comando, e a dificuldade no disparo e
                        coleta de resultados. Apesar disso, a ferramenta é a mais completa, no quesito geração de testes em redes MQTT por possuir a maior
                        quantidade de parâmetros editáveis, entre todas as estudadas. Assim, esse projeto busca facilitar a criação de experimentos com o MQTTLoader
                        e a coleta de resultados, oferecendo um sistema de coleta, poller, e uma interface gráfica que permitam um melhor gerenciamento dos experimentos.
                    </Text>
                </Article>
            </Container>
            <Container>
                <Article title={ {text: 'Sobre a Instituição'} }>

                </Article>
            </Container>
            <Container>
                <Article title={ {text: 'Como Contribuir'} }>
                    
                </Article>
            </Container>
        </>
    );
};

export default Homepage;