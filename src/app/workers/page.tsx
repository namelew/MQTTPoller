'use client'
import { Table, Thead, Tbody, Tr, Th, Td, Checkbox, Flex, Button } from "@chakra-ui/react";
import { useState } from "react";
import { IWorker } from "../interfaces/IWorker";

const Homepage = () => {
    const [selectAll, setSelectAll] = useState(false);
    const [selectedWorkers, setSelectedWorkers] = useState<string[]>([]);
    const [workers, setWorkers] = useState<IWorker[]>();

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

    const handleSelectRow = (id:string) => {
        if (selectedWorkers.includes(id)) {
            setSelectedWorkers(selectedWorkers.filter((rowId) => rowId !== id));
        } else {
            setSelectedWorkers([...selectedWorkers, id]);
        }
    };

    const startExperiment = () => {

    };

    return (
        <>
            <Flex justifyContent={'flex-end'}>
                <Button onClick={() => startExperiment()}>Iniciar Experimento</Button>
            </Flex>
            <Table variant="simple">
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
                    {workers?.map((worker) => (
                    <Tr key={worker.id}>
                        <Td>
                        <Checkbox
                            isChecked={selectedWorkers.includes(worker.id)}
                            onChange={() => handleSelectRow(worker.id)}
                        />
                        </Td>
                        <Td>{worker.id}</Td>
                        <Td>{worker.online ? 'Online' : 'Offline'}</Td>
                    </Tr>
                    ))}
                </Tbody>
            </Table>
        </>
    );
};

export default Homepage;