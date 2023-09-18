'use client'
import { Box, Link, Stack, Text } from "@chakra-ui/react";

function Footer() {
  return (
    <Box as="footer" role="contentinfo" mx="auto" maxW="7xl" py="12" px={{ base: '4', md: '8' }}>
      <Stack>
        <Text fontSize="sm">
          @Mais Informações e Contatos
        </Text>
        <Stack direction="row" spacing="4" align="center" justify="space-between">
          <Link href="https://github.com/namelew/MQTTPoller" rel="noopener noreferrer" target="_blank">
            Licença
          </Link>
          <Link href="https://github.com/namelew/MQTTPoller" rel="noopener noreferrer" target="_blank">
            Repositório
          </Link>
          <Link href="mailto:diogomaciel.cunha@gmail.com" rel="noopener noreferrer" target="_blank">
            Email
          </Link>
          <Link href="https://www.uffs.edu.br" rel="noopener noreferrer" target="_blank">
            Instituição
          </Link>
          <Link href="https://cc.uffs.edu.br/" rel="noopener noreferrer" target="_blank">
            Colegiado
          </Link>
        </Stack>
      </Stack>
    </Box>
  );
}

export default Footer;