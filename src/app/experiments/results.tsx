'use client'
import { IResult } from "interfaces/IExperiment";
import { saveAs } from 'file-saver';
import { Flex, Text, Button, HStack, Box, VStack } from "@chakra-ui/react";

interface Props {
    workerID: string
    result: IResult
}

const ResultContainer = ( { workerID, result } : Props) => {
    const onDownload = () => {
        const blob = new Blob([JSON.stringify(result)], {type: 'application/json'});
        saveAs(blob, `results-worker-${workerID}.json`);
    };

    const onLogFileDownload = () => {
        if (result.meta.log_file.name === '') {
            return
        }
        const blob = new Blob([result.meta.log_file.data], {type: 'text/csv'});
        saveAs(blob, `${result.meta.log_file.name}-${workerID}.${result.meta.log_file.extension}`);
    };

    return (
        <VStack key={workerID} p={5}>
            <HStack>
                <Text fontWeight="bold">Worker:</Text>
                <Text>{workerID}</Text>
            </HStack>
            <Box>
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
                    <Text>{result.subscribe.avg_latency.toString()} ms</Text>
                </Flex>
                <Flex gap='1'>
                    <Text fontWeight="bold">Assinatura Latência Máxima:</Text>
                    <Text>{result.subscribe.latency.toString()} ms</Text>
                </Flex>
                <Flex gap='1'>
                    <Text fontWeight="bold">Mensagens Recebidas:</Text>
                    <Text>{result.subscribe.received_messages}</Text>
                </Flex>
            </Box>
            <HStack justifyContent="space-between" width="100%">
                <Button
                    colorScheme='blue'
                    disabled={result.meta.log_file.name === ''}
                    onClick={onLogFileDownload}
                >Arquivo de Log</Button>
                <Button colorScheme='green' onClick={onDownload}>Baixar</Button>
            </HStack>
        </VStack>
    );
};

export default ResultContainer;