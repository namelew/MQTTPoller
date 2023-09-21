'use client'
import { Button, Flex, Checkbox, Table, Tr, Th, Thead, Tbody } from "@chakra-ui/react";
import { useState } from "react";
import { IWorker } from "interfaces/IWorker";
import Workers from "./data";
import ExperimentModal from "./modal";
import { useQuery } from "@tanstack/react-query";

interface Props {
    workers?:IWorker[],
}

const WorkersTable = ( { workers } : Props) => {
    const [selectAll, setSelectAll] = useState(false);
    const [openExperimentModal, setOpenExperimentModal] = useState(false);
    const [fetchError, setFetchError] = useState<Error>();
    const [selectedWorkers, setSelectedWorkers] = useState<string[]>([]);

    const { data, isLoading } = useQuery<IWorker[]>({
        queryKey: ['worker'],
        queryFn: () => fetch('/api/worker').then(res => res.json()).catch((error) => {
            setFetchError(error);
            return [];
        }),
        retry: 3,
        refetchInterval: 60000,
        initialData: workers,
    });

    if (!isLoading && data && !fetchError && data !== workers) {
        workers = data
    }

    if (fetchError) {
        console.log(fetchError);
    }

    const handleSelectAll = () => {
        if (workers) {
            if (selectAll) {
                setSelectedWorkers([]);
            } else {
                setSelectedWorkers(workers.map(worker => worker.id));
            }
            setSelectAll(!selectAll);
        }
    };

    const handlerOpenModal = () => {
        if (selectedWorkers.length > 0) {
            setOpenExperimentModal(true);
        }
    };

    const handlerCloseModal = () => {
        setOpenExperimentModal(false);
    };

    return (
        <>
            <Flex justifyContent={'flex-end'}>
                <Button 
                    onClick={handlerOpenModal}
                    disabled={selectedWorkers.length > 0}
                >Iniciar Experimento</Button>
                <ExperimentModal 
                    openModal={openExperimentModal}
                    onClose={handlerCloseModal}
                    selected={selectedWorkers}
                />
            </Flex>
            <Table variant='simple'>
                <Thead>
                    <Tr>
                    <Th>
                        <Checkbox isChecked={selectAll} onChange={handleSelectAll} />
                    </Th>
                    <Th>ID</Th>
                    <Th>Status</Th>
                    </Tr>
                </Thead>
                <Tbody>
                    <Workers 
                        workers={workers}
                        selectWorkers={setSelectedWorkers}
                        selectedWorkers={selectedWorkers}
                    />
                </Tbody>
            </Table>
        </>
    );
};

export default WorkersTable;