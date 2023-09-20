'use client'

import { Input, FormLabel, Box} from "@chakra-ui/react";

interface Props {
    label?: string
    name?: string,
    value?: number,
    onChange?: (event: React.FormEvent<HTMLInputElement>) => void
}

const Number = ( { label, value, name, onChange } : Props) => {
    return (
        <Box>
            {label && <FormLabel>{label}</FormLabel>}
            <Input type='number' value={value} name={name} onChange={onChange} />
        </Box>
    );
};

export default Number;