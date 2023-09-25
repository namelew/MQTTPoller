import { IResult } from "interfaces/IExperiment";
import { Flex, Text, Container } from "@chakra-ui/react";

interface Props {
    workerID: string
    result: IResult
}

const ResultContainer = ( { workerID, result } : Props) => {
    return (
        <Container key={workerID}>
            <Text fontWeight="bold">Worker:</Text>
            <Text>{workerID}</Text>
            <Flex gap='1'>
                <Text fontWeight="bold">Publicação Vazão Média:</Text>
                <Text>{result.publish.avg_throughput.toString()}</Text>
            </Flex>
            <Flex gap='1'>
                <Text fontWeight="bold">Publicação Vazão Máxima:</Text>
                <Text>{result.publish.max_throughput.toString()}</Text>
            </Flex>
            <Flex gap='1'>
                <Text fontWeight="bold">Publicação Vazão (Por segundo):</Text>
                <Text>Não implementado</Text>
            </Flex>
            <Flex gap='1'>
                <Text fontWeight="bold">Mensagens Publicadas:</Text>
                <Text>{result.publish.publiqued_messages}</Text>
            </Flex>
            <Flex gap='1'>
                <Text fontWeight="bold">Assinatura Vazão Média:</Text>
                <Text>{result.subscribe.avg_throughput.toString()}</Text>
            </Flex>
            <Flex gap='1'>
                <Text fontWeight="bold">Assinatura Vazão Máxima:</Text>
                <Text>{result.subscribe.max_throughput.toString()}</Text>
            </Flex>
            <Flex gap='1'>
                <Text fontWeight="bold">Assinatura Latência Média:</Text>
                <Text>{result.subscribe.avg_latency.toString()}</Text>
            </Flex>
            <Flex gap='1'>
                <Text fontWeight="bold">Assinatura Latência Máxima:</Text>
                <Text>{result.subscribe.latency.toString()}</Text>
            </Flex>
            <Flex gap='1'>
                <Text fontWeight="bold">Mensagens Recebidas:</Text>
                <Text>{result.subscribe.received_messages}</Text>
            </Flex>
        </Container>
    );
};

export default ResultContainer;