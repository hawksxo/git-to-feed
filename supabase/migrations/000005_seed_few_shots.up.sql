INSERT INTO few_shot_examples (uuid, archetype, input_context, expected_output, is_active)
VALUES
('seed-few-shot-release-01', 'RELEASE', 'Release v3.8-Omni-Flash: Major release with multi-modal audio-video 
understanding, agentic tools, and 89% cost reduction.', '🚀 ¡Presentamos la nueva versión v3.8-Omni-Flash de 
nuestro sistema!

Capacidades nativas de procesamiento multimodal, razonamiento largo y orquestación de herramientas reunidos en 
un solo lanzamiento.

Aspectos destacados:
- Inteligencia multimodal: Razonamiento sobre audio y video para workflows complejos y edición automatizada.
- Reducción del 89% en costos de procesamiento en comparación con la versión anterior.
- Ventana de contexto de 1M tokens con percepción precisa de eventos en tiempo real.

Revisa la documentación y release oficial aquí: https://github.com/hawksxo/git-to-feed/releases', true),

('seed-few-shot-feature-01', 'FEATURE', 'Feature: Added interactive Skills/Commands framework for 
automated repetitive workflows.', '✨ Incorporamos un nuevo sistema de Skills interactivas para automatizar tus tareas 
más repetitivas.

Piensa en una Skill como un atajo reutilizable para procesos recurrentes que funciona similar a un slash 
command. En lugar de configurar manualmente el contexto en cada petición, guardas las instrucciones una sola 
vez.

¿Cómo funciona?
Simplemente ejecuta la instrucción asignada y el sistema procesa el flujo completo por ti de manera 
transparente.

Revisa el Pull Request y la implementación en GitHub: https://github.com/hawksxo/git-to-feed/pull/38', true),

('seed-few-shot-refactor-01', 'REFACTOR', 'Refactor: Optimized execution engine, reduced token consumption by 
30% and improved latency.', '🛠️ Refactorización y optimización de rendimiento en el motor central.

Logramos una mejora sustancial en la velocidad de respuesta, reduciendo en más de un 30% el tiempo de 
procesamiento y el consumo de recursos en ejecuciones pesadas.

Detalles del cambio:
- Reestructuración de llamadas concurrentes para evitar asignaciones innecesarias en memoria.
- Menor uso de recursos logrando exactamente el mismo resultado y consistencia.

Revisa los detalles técnicos de la refactorización: https://github.com/hawksxo/git-to-feed', true)
ON CONFLICT (uuid) DO NOTHING;