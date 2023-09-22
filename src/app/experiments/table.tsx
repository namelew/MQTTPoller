'use client'
import { Flex, Table, Tr, Th, Thead, Tbody } from "@chakra-ui/react";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { IExperiment } from "interfaces/IExperiment";
import Experiments from "./data";

interface Props {
    experiments?:IExperiment[],
}

const ExperimentsTable = ( { experiments } : Props) => {
    const [fetchError, setFetchError] = useState<Error>();

    const { data, isLoading, error } = useQuery<IExperiment[]>({
        queryKey: ['experiment'],
        queryFn: () => fetch('/api/experiment').then(res => res.json()).catch((error) => {
            setFetchError(error);
            return [];
        }),
        retry: 3,
        refetchInterval: 60000,
        initialData: experiments,
    });

    if (!isLoading && data && !fetchError && data !== experiments) {
        experiments = data
    }

    if (fetchError || error) {
        console.log(fetchError, error);
    }

    return (
        <>
            <Flex justifyContent={'flex-end'}>
            </Flex>
            <Table variant='simple'>
                <Thead>
                    <Tr>
                        <Th>ID</Th>
                        <Th>Broker</Th>
                        <Th>Versão do Protocolo</Th>
                        <Th>Tópico</Th>
                        <Th>Tempo de Execução</Th>
                        <Th>Finalizado</Th>
                        <Th>Com Error?</Th>
                        <Th>Detalhes</Th>
                        <Th>Excluir</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    <Experiments experiments={experiments}/>
                </Tbody>
            </Table>
        </>
    );
};

export default ExperimentsTable;