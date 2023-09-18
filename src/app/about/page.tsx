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
                    <Text margin={1}>
                        A Universidade Federal da Fronteira Sul (UFFS) é uma instituição de ensino superior pública, popular e de qualidade. Criada pela Lei Nº 12.029,
                        de 15 de setembro de 2009, a UFFS abrange mais de 400 municípios da Mesorregião Grande Fronteira do Mercosul – Sudoeste do Paraná, Oeste de
                        Santa Catarina e Noroeste do Rio Grande do Sul.
                    </Text>
                    <Text margin={1}>
                        Contando com mais de 50 cursos de graduação, a Universidade já ultrapassou a marca de 8 mil alunos e completou, em 2022, treze anos de história.
                        As graduações oferecidas privilegiam as vocações da economia regional e estão em consonância com a Política Nacional de Formação de Professores
                        do Ministério da Educação (MEC). Para ingressar na UFFS é preciso realizar o ENEM, pois a Universidade atualmente adota o SiSU como método de acesso à graduação.
                    </Text>
                </Article>
            </Container>
        </>
    );
};

export default Homepage;