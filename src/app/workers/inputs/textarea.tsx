'use client'

import { Input, FormLabel, Box} from "@chakra-ui/react";

interface Props {
    label?: string
    name?: string,
    value?: string,
    type?: 'url' | 'text',
    onChange?: (event: React.FormEvent<any>) => void
}

const TextArea = ( { label, type = 'text', value ,name, onChange } : Props) => {
    return (
        <Box>
            {label && <FormLabel>{label}</FormLabel>}
            <Input type={type} value={value} name={name} onChange={onChange} />
        </Box>
    );
};

export default TextArea;